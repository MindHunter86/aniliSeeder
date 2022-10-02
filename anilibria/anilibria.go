package anilibria

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
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
	baseUrl      *url.URL
	loginUrl     *url.URL
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		gLog.Error().Err(err).Msg("there is some problems with cookiejar initialization because of internal golang error")
	}

	var apiClient *ApiClient = &ApiClient{
		http: &http.Client{
			Timeout:   time.Duration(gCli.Int("http-client-timeout")) * time.Second,
			Transport: httpTransport,
			Jar:       jar,
		},
	}

	// ??
	// todo optimize
	if err = apiClient.getLoginUrl(); err != nil {
		gLog.Error().Err(err).Msg("there are some errors in parsing login url; sleeping for 30 seconds")
		time.Sleep(30 * time.Second)
	} else {
		apiClient.checkAuthData()
	}

	return apiClient, apiClient.getBaseUrl()
}

func (m *ApiClient) getBaseUrl() (e error) {
	m.baseUrl, e = url.Parse(gCli.String("anilibria-api-baseurl"))
	return e
}

func (m *ApiClient) getLoginUrl() (e error) {
	m.loginUrl, e = url.Parse(gCli.String("anilibria-login-url"))
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
