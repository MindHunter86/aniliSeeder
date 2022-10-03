package anilibria

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type (
	rspGetSchedule struct {
		Day  int
		List []*rspGetTitle
	}
	rspGetTitle struct {
		Names    *rspTitleNames
		Status   *rspTitleStatus
		Type     *rspTitleType
		Torrents *rspTitleTorrents
	}
	rspTitleNames struct {
		Ru          string
		En          string
		Alternative string
	}
	rspTitleStatus struct {
		String string
		Code   int
	}
	rspTitleType struct {
		FullString string
		Code       int
		String     string
		Series     interface{}
		Length     int
	}
	rspTitleTorrents struct {
		Series *rspTorrentSeries
		List   []*rspTorrentList
	}
	rspTorrentList struct {
		TorrentId         int
		Series            *rspTorrentSeries
		Quality           *rspTorrentQuality
		Leechers          int
		Seeders           int
		Downloads         int
		TotalSize         int64
		Url               string
		UploadedTimestamp *time.Time
		Hash              string
		Metadata          interface{}
		RawBase64File     interface{}
	}
	rspTorrentSeries struct {
		Firest int
		Last   int
		String string
	}
	rspTorrentQuality struct {
		String     string
		Type       string
		Resolution string
		Encoder    string
		LqAudio    interface{}
	}
)

type ApiRequestMethod string

const (
	apiMethodGetSchedule ApiRequestMethod = "/getSchedule"
)

type SiteRequestMethod string

const (
	siteMethodLogin           SiteRequestMethod = "/public/login.php"
	siteMethodTorrentDownload SiteRequestMethod = "/public/torrent/download.php"
)

var (
	errApiAuthorizationFailed = errors.New("there is some problems with the authorization proccess")
	errApiAbnormalResponse    = errors.New("there is some problems with anilibria servers communication")
)

// common
func (*ApiClient) debugHttpHandshake(data interface{}) {
	if !gCli.Bool("debug") {
		return
	}

	var dump []byte

	switch v := data.(type) {
	case *http.Request:
		dump, _ = httputil.DumpRequestOut(data.(*http.Request), false)
	case *http.Response:
		dump, _ = httputil.DumpResponse(data.(*http.Response), false)
	default:
		gLog.Error().Msgf("there is an internal application error; undefined type - %T", v)
	}

	log.Println(string(dump))
}

func (m *ApiClient) apiAuthorize(authBody io.Reader) (e error) {
	gLog.Debug().Msg("Called apiAuthorize")

	var loginUrl = m.siteBaseUrl.String() + string(siteMethodLogin)

	var req *http.Request
	if req, e = http.NewRequest("POST", loginUrl, authBody); e != nil {
		return
	}

	req = m.getBaseRequest(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}
	defer rsp.Body.Close()

	var jar *cookiejar.Jar
	if jar, e = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); e != nil {
		gLog.Error().Err(e).Msg("there is some problems with cookiejar initialization because of internal golang error")
	}

	jar.SetCookies(m.siteBaseUrl, rsp.Cookies())
	m.http.Jar = jar

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		cookies := m.http.Jar.Cookies(m.siteBaseUrl)
		for _, cookie := range cookies {
			if cookie.Name == "PHPSESSID" && cookie.Value != "" {
				gLog.Info().Str("session", cookie.Value).Msg("authentication proccess has been completed successfully")
				return nil
			}
		}

		gLog.Error().Int("login_response_code", rsp.StatusCode).Msg("there is abnormal site reponse; auth failed on 200 OK; check logs")
		return errApiAuthorizationFailed
	default:
		gLog.Error().Int("login_response_code", rsp.StatusCode).Msg("there is abnormal status code from login page; check you auth data")
		return errApiAuthorizationFailed
	}

	return
}

func (*ApiClient) getBaseRequest(req *http.Request) *http.Request {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en,ru;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("DNT", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	return req
}

func (m *ApiClient) checkApiAuthorization(rrl *url.URL) error {
	if m.unauthorized {
		gLog.Debug().Msg("authorization step has been skipped because of `unauthorized` flag detected")
		return nil
	}

	if m.http.Jar == nil || len(m.http.Jar.Cookies(rrl)) == 0 {
		gLog.Info().Msg("unauthorized request detected; initiate the authentication proccess...")
		return m.GetApiAuthorization()
	}

	for _, cookie := range m.http.Jar.Cookies(rrl) {
		if cookie.Name == "PHPSESSID" && cookie.Value != "" {
			return nil
		}
	}

	gLog.Warn().Msg("there is no PHPSESSID cookie found; initiate the authentication proccess...")
	return m.GetApiAuthorization()
}

func (m *ApiClient) getTorrentFile(titleId string) (e error) {
	var rrl *url.URL
	if rrl, e = url.Parse(m.siteBaseUrl.String() + string(siteMethodTorrentDownload)); e != nil {
		return
	}

	var rgs = &url.Values{}
	rgs.Add("id", titleId)
	rrl.RawQuery = rgs.Encode()

	if e = m.checkApiAuthorization(rrl); e != nil {
		return
	}

	var req *http.Request
	if req, e = http.NewRequest("GET", rrl.String(), nil); e != nil {
		return
	}
	req = m.getBaseRequest(req) // ???

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}
	defer rsp.Body.Close()

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		gLog.Debug().Msg("the requested torrent file has been found")
	default:
		gLog.Warn().Int("response_code", rsp.StatusCode).Msg("could not fetch the requested torrent file because of abnormal anilibria server response")
		return errApiAbnormalResponse
	}

	if rsp.Header.Get("Content-Type") != "application/x-bittorrent" {
		gLog.Warn().Msg("there is an abnormal content-type in the torrent file response")
	}

	_, params, e := mime.ParseMediaType(rsp.Header.Get("Content-Disposition"))
	if e != nil {
		return
	}

	gLog.Debug().Str("filename", params["filename"]).Msg("trying to download and save the torrent file...")
	return m.parseFileFromResponse(&rsp.Body, params["filename"])
}

func (*ApiClient) parseFileFromResponse(rsp *io.ReadCloser, filename string) (e error) {

	var fi fs.FileInfo
	if fi, e = os.Stat(gCli.String("torrentfiles-dir") + "/" + filename); e != nil {
		if !os.IsNotExist(e) {
			return
		}

		gLog.Debug().Str("path", gCli.String("torrentfiles-dir")+"/"+filename).
			Msg("given filepath is not found; continue...")
	}

	if fi != nil {
		gLog.Warn().Str("path", gCli.String("torrentfiles-dir")+"/"+filename).
			Msg("destination file isa already exists; skipping downloading...")
		return
	}

	var fd *os.File
	if fd, e = os.Create(gCli.String("torrentfiles-dir") + "/" + filename); e != nil {
		return
	}
	defer fd.Close()

	var n int64
	if n, e = io.Copy(fd, *rsp); e != nil {
		return
	} else {
		gLog.Info().Int64("bytes", n).Msg("the torrnet file has been successfully saved")
		return
	}
}

func (m *ApiClient) getApiResponse(httpMethod string, apiMethod ApiRequestMethod, rspSchema interface{}) (e error) {
	gLog.Debug().Msg("Called getResponse.")

	var rrl = *m.apiBaseUrl
	rrl.Path = rrl.Path + string(apiMethod)

	if e = m.checkApiAuthorization(&rrl); e != nil {
		return
	}

	var req *http.Request
	if req, e = http.NewRequest(httpMethod, rrl.String(), nil); e != nil {
		return
	}

	req = m.getBaseRequest(req) // ???

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}
	defer rsp.Body.Close()

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		gLog.Info().Str("api_method", string(apiMethod)).Msg("Correct response")
	default:
		gLog.Warn().Str("api_method", string(apiMethod)).Int("api_response_code", rsp.StatusCode).Msg("Abnormal API response")
		return errApiAbnormalResponse
	}

	return m.parseResponse(&rsp.Body, rspSchema)
}

func (*ApiClient) parseResponse(rsp *io.ReadCloser, schema interface{}) error {
	if data, err := ioutil.ReadAll(*rsp); err == nil {
		return json.Unmarshal(data, &schema)
	} else {
		return err
	}
}

// methods
func (m *ApiClient) GetApiAuthorization() (e error) {
	gLog.Debug().Msg("Called apiAuthorize()")

	authForm := url.Values{
		"mail":    {gCli.String("anilibria-login-username")},
		"passwd":  {gCli.String("anilibria-login-password")},
		"fa2code": {""},
		"csrf":    {"1"},
	}

	gLog.Debug().Str("username", gCli.String("anilibria-login-username")).Msg("trying to complete authentication proccess on anilibria")
	return m.apiAuthorize(strings.NewReader(authForm.Encode()))
}

func (m *ApiClient) GetTitleSchedule() (schedule []*rspGetSchedule, e error) {
	if e = m.getApiResponse("GET", apiMethodGetSchedule, &schedule); e != nil {
		return
	}

	gLog.Info().Int("response_length", len(schedule)).Msg("DONE!")
	return
}

func (m *ApiClient) getTitleTorrentFile(torrentId string) (e error) {
	gLog.Debug().Msg("trying to fetch torrent file for " + torrentId)
	return m.getTorrentFile(torrentId)
}
