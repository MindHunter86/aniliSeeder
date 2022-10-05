package swarm

import (
	"context"

	"github.com/MindHunter86/aniliSeeder/anilibria"
	"github.com/MindHunter86/aniliSeeder/deluge"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	gCli    *cli.Context
	gLog    *zerolog.Logger
	gCtx    context.Context
	gDeluge *deluge.Client
	gAniApi *anilibria.ApiClient
)

type Swarm interface {
	Bootstrap() error
}
