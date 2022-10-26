package anilibria

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
)

var (
	gCli *cli.Context
	gLog *zerolog.Logger
)

type ApiClient struct {
	http         *http.Client
	apiBaseUrl   *url.URL
	siteBaseUrl  *url.URL
	unauthorized bool
}

// TODO:
// - Check keepalive!!

func NewApiClient(ctx *cli.Context, log *zerolog.Logger) (*ApiClient, error) {
	gCli, gLog = ctx, log

	defaultTransportDialContext := func(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
		return dialer.DialContext
	}

	http1Transport := &http.Transport{
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   gCli.Duration("http-tcp-timeout"),
			KeepAlive: gCli.Duration("http-keepalive-timeout"),
		}),

		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: gCli.Bool("http-client-insecure"),
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		},
		TLSHandshakeTimeout: gCli.Duration("http-tls-handshake-timeout"),

		MaxIdleConns:    gCli.Int("http-max-idle-conns"),
		IdleConnTimeout: gCli.Duration("http-idle-timeout"),

		DisableCompression: false,
		DisableKeepAlives:  false,
		ForceAttemptHTTP2:  true,
	}

	var httpTransport http.RoundTripper = http1Transport
	http2Transport, err := http2.ConfigureTransports(http1Transport)
	if err != nil {
		httpTransport = http2Transport
		gLog.Warn().Err(err).Msg("could not upgrade http transport to v2 because of internal error")
	}

	var apiClient = &ApiClient{
		http: &http.Client{
			Timeout:   time.Duration(gCli.Int("http-client-timeout")) * time.Second,
			Transport: httpTransport,
		},
	}

	// ??
	// todo optimize
	if err = apiClient.getSiteBaseUrl(); err != nil {
		gLog.Error().Err(err).Msg("there are some errors in parsing login url; sleeping for 30 seconds")
		time.Sleep(30 * time.Second)
	} else {
		apiClient.checkAuthData()
	}

	return apiClient, apiClient.getApiBaseUrl()
}

func (m *ApiClient) getApiBaseUrl() (e error) {
	m.apiBaseUrl, e = url.Parse(gCli.String("anilibria-api-baseurl"))
	return e
}

func (m *ApiClient) getSiteBaseUrl() (e error) {
	m.siteBaseUrl, e = url.Parse(gCli.String("anilibria-baseurl"))
	return e
}

func (m *ApiClient) checkAuthData() {
	if gCli.String("anilibria-login-username") == "" || gCli.String("anilibria-login-password") == "" {
		m.unauthorized = true
		gLog.Warn().Msg("could not parse username and\\or password")
		gLog.Warn().Msg("ATTENTION! Unauthorized peering detected; Anilibria could not detect the seeder, so it's stats will be disable!!!")
		gLog.Info().Msg("\"unauthorized\" has been toggled; sleeping for 3 seconds...")
		time.Sleep(3 * time.Second)
	}
}

// popular domains origin
// - getSchedule
// - online top
// - releases // ??

// https://api.anilibria.tv/v2/getSchedule?days=0&filter=id,code,names,updated,last_change,status,type,torrents

func (*ApiClient) GetPopularTorrents() (error, error) {
	return nil, nil
}

func (m *TitleTorrent) GetShortHash() string {
	return m.Hash[0:9]
}
