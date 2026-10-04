package app

import (
	"ITK_Code/m/v2/internal/config"
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"
)

func Run(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {

		return fmt.Errorf(
			"error loading config file path: %s, error: %s",
			cfgPath,
			err,
		)
	}

	secret, err := os.ReadFile(cfg.JWTSecretPath)
	if err != nil {

		return fmt.Errorf(
			"error reading JWT secret from %s: %s",
			cfg.JWTSecretPath,
			err,
		)
	}

	fxApp := fx.New(
		fx.Supply(cfg, string(secret)),
		fx.Provide(New),
		fx.Invoke(func(lifecycle fx.Lifecycle, application *App) {
			lifecycle.Append(fx.Hook{
				OnStart: func(context.Context) error {
					application.Start()
					return nil
				},
				OnStop: func(context.Context) error {
					application.Stop()
					return nil
				},
			})
		}),
	)
	if err := fxApp.Err(); err != nil {
		return fmt.Errorf("create application failed: %w", err)
	}
	fxApp.Run()

	return nil
}
