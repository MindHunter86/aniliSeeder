package anilibria

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type (
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

// common
func (m *ApiClient) checkApiHealth() {
	return
}

func (m *ApiClient) getBaseRequest(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en,ru;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("DNT", "1")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Sec-GPC", "1")
	// req.Header.Set("Pragma", "no-cache")
	// req.Header.Set("Cache-Control", "no-cache")
}

func (m *ApiClient) getResponse(httpMethod string, apiMethod ApiRequestMethod, rspSchema interface{}) (e error) {
	gLog.Debug().Msg("Called getResponse.")

	var reqUrl url.URL = *m.baseUrl
	reqUrl.Path = reqUrl.Path + string(apiMethod)

	var req *http.Request
	if req, e = http.NewRequest(httpMethod, reqUrl.String(), nil); e != nil {
		return
	}

	m.getBaseRequest(req) // ???

	if gCli.Bool("debug") {
		var dump []byte
		dump, e = httputil.DumpRequestOut(req, true)
		if e != nil {
			gLog.Warn().Err(e).Msg("could not dump the request because of httputil internal errors")
		}

		// gLog.Debug().Bytes("request", dump).Msg("")
		log.Println(string(dump))
	}

	var rsp *http.Response
	if rsp, e = m.http.Do(req); e != nil {
		return
	}

	if rsp.Uncompressed {
		gLog.Warn().Str("api_method", string(apiMethod)).Msg("there is uncopressed request detected")
	}

	if gCli.Bool("debug") {
		var dump []byte
		dump, e = httputil.DumpResponse(rsp, false)
		if e != nil {
			gLog.Warn().Err(e).Msg("could not dump the request because of httputil internal errors")
		}

		// gLog.Debug().Bytes("response", dump).Msg("")
		log.Println(string(dump))
	}

	switch rsp.StatusCode {
	case http.StatusOK:
		gLog.Info().Str("api_method", string(apiMethod)).Msg("Correct response")
	default:
		gLog.Warn().Str("api_method", string(apiMethod)).Int("api_response_code", rsp.StatusCode).Msg("Abnormal API response")
	}

	defer rsp.Body.Close()
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
