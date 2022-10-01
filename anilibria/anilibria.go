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
	http    *http.Client
	baseUrl *url.URL
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

	var apiClient *ApiClient = &ApiClient{
		http: &http.Client{
			Timeout:   time.Duration(gCli.Int("http-client-timeout")) * time.Second,
			Transport: httpTransport,
		},
	}

	return apiClient, apiClient.getBaseUrl()
}

func (m *ApiClient) getBaseUrl() (e error) {
	m.baseUrl, e = url.Parse(gCli.String("anilibria-api-baseurl"))
	return e
}
