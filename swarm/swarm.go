package swarm

import (
	"context"

	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	gCli    *cli.Context
	gLog    *zerolog.Logger
	gCtx    context.Context
	gDeluge *deluge.Client
)

var (
	swarmCA []byte = []byte(`-----BEGIN CERTIFICATE-----
MIICGDCCAb+gAwIBAgICB+MwCgYIKoZIzj0EAwIwdTELMAkGA1UEBhMCVVMxCTAH
BgNVBAgTADEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEbMBkGA1UECRMSR29sZGVu
IEdhdGUgQnJpZGdlMQ4wDAYDVQQREwU5NDAxNjEWMBQGA1UEChMNQ29tcGFueSwg
SU5DLjAeFw0yMjEwMDUxNjI5NDhaFw0zMjEwMDUxNjI5NDhaMHUxCzAJBgNVBAYT
AlVTMQkwBwYDVQQIEwAxFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xGzAZBgNVBAkT
EkdvbGRlbiBHYXRlIEJyaWRnZTEOMAwGA1UEERMFOTQwMTYxFjAUBgNVBAoTDUNv
bXBhbnksIElOQy4wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARZWeW7xifJLyLY
jhhbhCxiYGxqmbFBfBka5gfNcQZnuZtAyeWv2ClHqREoQLWq5wVzgf2vE3jj7sVB
plJOTW4Toz8wPTAOBgNVHQ8BAf8EBAMCAoQwHQYDVR0lBBYwFAYIKwYBBQUHAwIG
CCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0EAwIDRwAwRAIgBeJBEtMf
uNxtx64bnu+f2U1mZi2zFVd0/QKBGfcs+LcCIAId/yKdpsYYbhvw/b4My2OGwNzd
m0OYysOoZWEzunDG
-----END CERTIFICATE-----`)
)

type Swarm interface {
	Bootstrap() error
}
