package main

import (
	"context"

	"go.uber.org/fx"

	appfx "github.com/your-org/boilerplate-go/internal/fx"
	"github.com/your-org/boilerplate-go/internal/server"
)

func main() {
	fx.New(
		appfx.AppModule,
		fx.Invoke(startServer),
	).Run()
}

func startServer(lc fx.Lifecycle, srv *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				srv.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
