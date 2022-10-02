package anilibria

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type (
	reqAuthForm struct {
		Mail    string
		Passwd  string
		Fa2Code string
		Csrf    int
	}
	rspGetSchedule struct {
		Day  int
		List []*rspGetTile
	}
	rspGetTile struct {
		Names    *rspTileNames
		Status   *rspTileStatus
		Type     *rspTileType
		Torrents *rspTileTorrents
	}
	rspTileNames struct {
		Ru          string
		En          string
		Alternative string
	}
	rspTileStatus struct {
		String string
		Code   int
	}
	rspTileType struct {
		FullString string
		Code       int
		String     string
		Series     interface{}
		Length     int
	}
	rspTileTorrents struct {
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

var (
	errApiAuthorizationFailed = errors.New("there is some problems with the authorization proccess")
)

// common
func (m *ApiClient) checkApiHealth() {
	return
}

func (m *ApiClient) debugHttpHandshake(data interface{}) {
	if !gCli.Bool("debug") {
		return
	}

	var e error
	var dump []byte

	switch v := data.(type) {
	case *http.Request:
		dump, e = httputil.DumpRequestOut(data.(*http.Request), true)
	case *http.Response:
		dump, e = httputil.DumpResponse(data.(*http.Response), false)
	default:
		gLog.Error().Msgf("there is an internal application error; undefined type - %T", v)
	}

	if e != nil {
		gLog.Warn().Msg("could not dump the request because of httputil internal errors")
		return
	}

	log.Println(string(dump))
}

func (m *ApiClient) apiAuthorize(authBody io.Reader) (e error) {
	gLog.Debug().Msg("Called apiAuthorize")

	var req *http.Request
	if req, e = http.NewRequest("POST", m.loginUrl.String(), authBody); e != nil {
		return
	}

	m.getBaseRequest(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}
	defer rsp.Body.Close()

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		cookies := m.http.Jar.Cookies(m.loginUrl)
		for _, cookie := range cookies {
			if cookie.Name == "PHPSESSID" {
				gLog.Info().Str("session", cookie.Value).Msg("authorization process has been completed successfully")
				gLog.Info().Time("session_expire", cookie.Expires).Msg("authorization process has been completed successfully")
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

func (m *ApiClient) getBaseRequest(req *http.Request) {
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
}

// ??
// todo
// refactor
func (m *ApiClient) checkApiAuthorization(reqUrl *url.URL) error {
	if m.unauthorized == true {
		gLog.Debug().Msg("authorization step has been skipped because of `unauthorized` flag detected")
		return nil
	}

	if len(m.http.Jar.Cookies(reqUrl)) == 0 {
		gLog.Info().Msg("unauthorized request detected; initiate the authorization proccess...")
		return m.GetApiAuthorization()
	}

	for _, cookie := range m.http.Jar.Cookies(reqUrl) {
		if cookie.Name == "PHPSESSID" {
			if time.Now().Unix() < cookie.Expires.Unix() {
				gLog.Info().Time("session_expire", cookie.Expires).Msg("session expiration has been verified")
				return nil
			} else {
				gLog.Warn().Time("now", time.Now()).Time("session_expire", cookie.Expires).Msg("session has been expired; initiating the reauthentication proccess")
				m.cleanApiAuthorization(reqUrl)
				return m.GetApiAuthorization()
			}

			gLog.Warn().Msg("there is no PHPSESSID cookie found; initiate the authentication proccess...")
			return m.GetApiAuthorization()
		}
	}

	return nil
}

func (m *ApiClient) cleanApiAuthorization(reqUrl *url.URL) {
	m.http.Jar.SetCookies(reqUrl, nil)
}

func (m *ApiClient) getResponse(httpMethod string, apiMethod ApiRequestMethod, rspSchema interface{}) (e error) {
	gLog.Debug().Msg("Called getResponse.")

	var reqUrl url.URL = *m.baseUrl
	reqUrl.Path = reqUrl.Path + string(apiMethod)

	if e = m.checkApiAuthorization(&reqUrl); e != nil {
		return
	}

	var req *http.Request
	if req, e = http.NewRequest(httpMethod, reqUrl.String(), nil); e != nil {
		return
	}

	m.getBaseRequest(req) // ???

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
	}

	return m.parseResponse(&rsp.Body, rspSchema)
}

func (m *ApiClient) parseResponse(rsp *io.ReadCloser, schema interface{}) error {
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

	gLog.Debug().Str("username", gCli.String("anilibria-login-username")).Msg("trying to complete authorization process on anilibria")
	return m.apiAuthorize(strings.NewReader(authForm.Encode()))
}

func (m *ApiClient) GetTileSchedule() (e error) {
	gLog.Debug().Msg("Called GetTileSchedule")

	var schedule []*rspGetSchedule

	if e = m.getResponse("GET", apiMethodGetSchedule, &schedule); e != nil {
		gLog.Debug().Msg("Called GetTileSchedule 2")
		return e
	}

	gLog.Debug().Msg("Called GetTileSchedule 3")

	gLog.Info().Int("response_length", len(schedule)).Msg("DONE!")
	return
}
