package anilibria

import (
	"io"
	"net/http"
)

type SiteRequestMethod string

const (
	siteMethodLogin           SiteRequestMethod = "/public/login.php"
	siteMethodTorrentDownload SiteRequestMethod = "/public/torrent/download.php"
	siteMethodSessions        SiteRequestMethod = "/pages/cp.php"
)

func (m *ApiClient) getSiteResponse(hmethod string, smethod SiteRequestMethod, payload ...interface{}) (_ *[]byte, e error) {
	rrl := *m.siteBaseUrl
	rrl.Path += string(smethod)

	if e = m.checkApiAuthorization(&rrl); e != nil {
		return
	}

	var body io.ReadCloser
	if len(payload) != 0 {
		body = payload[0].(io.ReadCloser)
	}

	var req *http.Request
	if req, e = http.NewRequest(hmethod, rrl.String(), body); e != nil {
		return
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

func (m *ApiClient) getSessionsPage() (*[]byte, error) {
	return m.getSiteResponse("GET", siteMethodSessions)
}
