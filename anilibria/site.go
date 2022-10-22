package anilibria

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SiteRequestMethod string

const (
	siteMethodLogin           SiteRequestMethod = "/public/login.php"
	siteMethodTorrentDownload SiteRequestMethod = "/public/torrent/download.php"
	siteMethodSessions        SiteRequestMethod = "/pages/cp.php"
	siteMiethodCloseSession   SiteRequestMethod = "/public/close.php"
)

func (m *ApiClient) getSiteResponse(hmethod string, smethod SiteRequestMethod, payload ...interface{}) (_ *[]byte, e error) {
	rrl := *m.siteBaseUrl
	rrl.Path += string(smethod)

	if e = m.checkApiAuthorization(&rrl); e != nil {
		return
	}

	var body io.Reader
	if len(payload) != 0 {
		body = payload[0].(io.Reader)
	}

	var req *http.Request
	if req, e = http.NewRequest(hmethod, rrl.String(), body); e != nil {
		return
	}

	if smethod == siteMiethodCloseSession {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}
	defer rsp.Body.Close()

	m.debugHttpHandshake(req)
	m.debugHttpHandshake(rsp)

	if rsp.StatusCode != http.StatusOK {
		gLog.Warn().Int("response_code", rsp.StatusCode).Msg("abormal response from the site")
	}

	var buf []byte
	if buf, e = io.ReadAll(rsp.Body); e != nil {
		return
	}

	return &buf, e
}

func (*ApiClient) parseJSONReponse(payload *[]byte, schema interface{}) error {
	return json.Unmarshal(*payload, &schema)
}

func (m *ApiClient) getSessionsPage() (*[]byte, error) {
	return m.getSiteResponse("GET", siteMethodSessions)
}

func (m *ApiClient) dropActiveSession(sid string) (_ bool, e error) {
	dropResponseChema := struct {
		Err string
		Mes string
		Key string `json:",omitempty"`
	}{}

	reqPayload := url.Values{
		"id":         {sid},
		"csrf_token": {""},
	}

	var rspPayload *[]byte
	if rspPayload, e = m.getSiteResponse("POST", siteMiethodCloseSession, strings.NewReader(reqPayload.Encode())); e != nil {
		return
	}

	if e = m.parseJSONReponse(rspPayload, &dropResponseChema); e != nil {
		return
	}

	gLog.Debug().Str("response_err", dropResponseChema.Err).Str("response_mes", dropResponseChema.Mes).
		Str("response_key", dropResponseChema.Key).Msg("")
	return dropResponseChema.Err == "ok" && dropResponseChema.Mes == "Success", e
}
