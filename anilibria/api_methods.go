package anilibria

import (
	"net/http"
	"net/url"
	"strings"
)

func (m *ApiClient) GetApiAuthorization() (e error) {
	gLog.Debug().Msg("Called apiAuthorize()")

	authForm := url.Values{
		"mail":   {gCli.String("anilibria-login-username")},
		"passwd": {gCli.String("anilibria-login-password")},
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

func (m *ApiClient) GetTitleTorrentFile(tid string) (string, *[]byte, error) {
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

// V2

func (m *ApiClient) SearchTitlesByName(name string) (titles []*Title, e error) {
	params := []string{"search", name}
	e = m.getApiResponseV2(http.MethodGet, apiMethodSearchTitles, params).parseApiResponse(&titles)
	return
}
