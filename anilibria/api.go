package anilibria

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
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
	apiError struct {
		Error *apiErrorDetails
	}
	apiErrorDetails struct {
		Code    int
		Message string
	}
	TitleSchedule struct {
		Day  int
		List []*Title
	}
	Title struct {
		Id         int
		Code       string
		Updated    uint64 // sometimes the anilibria project mark their update time as a NULL
		LastChange uint64 `json:"last_change"` // I dont know how to mark this fields as "if time.Parse fails - ignore"
		Names      *TitleNames
		Status     *TitleStatus
		Type       *TitleType
		Torrents   *TitleTorrents
	}
	TitleNames struct {
		Ru          string
		En          string
		Alternative string
	}
	TitleStatus struct {
		String string
		Code   int
	}
	TitleType struct {
		FullString string `json:"full_string"`
		Code       int
		String     string
		Series     interface{}
		Length     int
	}
	TitleTorrents struct {
		Series *TorrentSeries
		List   []*TitleTorrent
	}
	TitleTorrent struct {
		TorrentId         int `json:"torrent_id"`
		Series            *TorrentSeries
		Quality           *TorrentQuality
		Leechers          int
		Seeders           int
		Downloads         int
		TotalSize         int64 `json:"total_size"`
		Url               string
		UploadedTimestamp *time.Time
		Hash              string
		Metadata          interface{}
		RawBase64File     interface{}
	}
	TorrentSeries struct {
		Firest int
		Last   int
		String string
	}
	TorrentQuality struct {
		String     string
		Type       string
		Resolution string
		Encoder    string
		LqAudio    interface{} `json:"lq_audio"`
	}
)

// https://api.anilibria.tv/v2/getSchedule?days=0&filter=id,code,names,updated,last_change,status,type,torrents
const defaultApiMethodFilter = "id,code,names,updated,last_change,status,type,torrents"
const defaultApiMethodLimit = "10"

type ApiRequestMethod string

const (
	apiMethodGetSchedule ApiRequestMethod = "/getSchedule"
	apiMethodGetUpdates  ApiRequestMethod = "/getUpdates"
	apiMethodGetChanges  ApiRequestMethod = "/getChanges"
)

type SiteRequestMethod string

const (
	siteMethodLogin           SiteRequestMethod = "/public/login.php"
	siteMethodTorrentDownload SiteRequestMethod = "/public/torrent/download.php"
)

var (
	errApiAuthorizationFailed = errors.New("there is some problems with the authorization process")
	errApiAbnormalResponse    = errors.New("there is some problems with anilibria servers communication")
)

// common
func (*ApiClient) debugHttpHandshake(data interface{}) {
	if !gCli.Bool("http-debug") {
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
				gLog.Info().Str("session", cookie.Value).Msg("authentication process has been completed successfully")
				return nil
			}
		}

		gLog.Error().Int("login_response_code", rsp.StatusCode).Msg("there is abnormal site reponse; auth failed on 200 OK; check logs")
		return errApiAuthorizationFailed
	default:
		gLog.Error().Int("login_response_code", rsp.StatusCode).Msg("there is abnormal status code from login page; check you auth data")
		return errApiAuthorizationFailed
	}
}

func (*ApiClient) getBaseRequest(req *http.Request) *http.Request {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en,ru;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("Connection", "keep-alive") // !!
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
		gLog.Info().Msg("unauthorized request detected; initiate the authentication process...")
		return m.GetApiAuthorization()
	}

	for _, cookie := range m.http.Jar.Cookies(rrl) {
		if cookie.Name == "PHPSESSID" && cookie.Value != "" {
			return nil
		}
	}

	gLog.Warn().Msg("there is no PHPSESSID cookie found; initiate the authentication process...")
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
			Msg("given filepath was not found; continue...")
	}

	if fi != nil {
		gLog.Warn().Str("path", gCli.String("torrentfiles-dir")+"/"+filename).
			Msg("destination file is already exists; skipping downloading...")
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
	}

	gLog.Info().Int64("bytes", n).Msg("the torrnet file has been successfully saved")
	return fd.Sync()
}

func (m *ApiClient) getApiResponse(httpMethod string, apiMethod ApiRequestMethod, rspSchema interface{}) (e error) {
	gLog.Debug().Msg("Called getResponse.")

	var rrl = *m.apiBaseUrl
	rrl.Path = rrl.Path + string(apiMethod)

	var rgs = &url.Values{}
	rgs.Add("filter", defaultApiMethodFilter)
	// rgs.Add("limit", defaultApiMethodLimit)
	rrl.RawQuery = rgs.Encode()

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
		gLog.Debug().Msg("trying to get error description")

		// apierr := &apiError{}
		var apierr *apiError
		if e = m.parseResponse(&rsp.Body, &apierr); e != nil {
			gLog.Error().Err(e).Msg("could not get api error description")
			return errApiAbnormalResponse
		}

		gLog.Warn().Int("error_code", apierr.Error.Code).Str("error_desc", apierr.Error.Message).Msg("api error has been parsed")
		return errApiAbnormalResponse
	}

	return m.parseResponse(&rsp.Body, rspSchema)
}

func (*ApiClient) parseResponse(rsp *io.ReadCloser, schema interface{}) error {
	if data, err := io.ReadAll(*rsp); err == nil {
		return json.Unmarshal(data, &schema)
	} else {
		return err
	}
}

func (m *ApiClient) downloadTorrentFile(tid string) (_ string, _ *io.ReadCloser, e error) {
	var rrl *url.URL
	if rrl, e = url.Parse(m.siteBaseUrl.String() + string(siteMethodTorrentDownload)); e != nil {
		return
	}

	var rgs = &url.Values{}
	rgs.Add("id", tid)
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
	// the caller must close the returned body
	// defer rsp.Body.Close()

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		gLog.Debug().Msg("the requested torrent file has been found")
	default:
		gLog.Warn().Int("response_code", rsp.StatusCode).Msg("could not fetch the requested torrent file because of abnormal anilibria server response")
		return "", nil, errApiAbnormalResponse
	}

	if rsp.Header.Get("Content-Type") != "application/x-bittorrent" {
		gLog.Warn().Msg("there is an abnormal content-type in the torrent file response")
	}

	_, params, e := mime.ParseMediaType(rsp.Header.Get("Content-Disposition"))
	if e != nil {
		return
	}

	gLog.Debug().Str("filename", params["filename"]).Msg("trying to download and save the torrent file...")
	return params["filename"], &rsp.Body, e
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

	gLog.Debug().Str("username", gCli.String("anilibria-login-username")).Msg("trying to complete authentication process on anilibria")
	return m.apiAuthorize(strings.NewReader(authForm.Encode()))
}

func (m *ApiClient) GetTitlesFromSchedule() (titles []*Title, e error) {
	var weekSchedule []*TitleSchedule
	if e = m.getApiResponse("GET", apiMethodGetSchedule, &weekSchedule); e != nil {
		return
	}

	for _, schedule := range weekSchedule {
		titles = append(titles, schedule.List...)
	}

	gLog.Debug().Int("titles_count", len(titles)).Msg("titles has been successfully parsed from schedule")
	return
}

func (m *ApiClient) SaveTitleTorrentFile(torrentId string) (e error) {
	gLog.Debug().Msg("trying to fetch torrent file for " + torrentId)
	return m.getTorrentFile(torrentId)
}

func (m *ApiClient) GetTitleTorrentFile(tid string) (string, *io.ReadCloser, error) {
	return m.downloadTorrentFile(tid)
}

func (m *ApiClient) GetLastUpdates() (titles []*Title, e error) {
	if e = m.getApiResponse("GET", apiMethodGetUpdates, &titles); e != nil {
		return
	}

	return titles, e
}

func (m *ApiClient) GetLastChanges() (titles []*Title, e error) {
	if e = m.getApiResponse("GET", apiMethodGetChanges, &titles); e != nil {
		return
	}

	return titles, e
}

func (m *ApiClient) GetTitlesSchedule() (schedule []*TitleSchedule, e error) {
	if e = m.getApiResponse("GET", apiMethodGetSchedule, &schedule); e != nil {
		return
	}

	return
}
