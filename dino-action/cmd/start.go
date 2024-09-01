package cmd

import (
	"context"
	"fmt"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/config"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/handlers"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/store"
	"gitlab.com/stackvista/demo/kubecon2024/poi/pkg/otel"
	"gitlab.com/stackvista/demo/kubecon2024/poi/pkg/reaper"
)

func NewStartCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the Dino Action service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context(), cfg)
		},
	}
}

func start(ctx context.Context, cfg *config.Config) error {
	if err := config.Validate(cfg); err != nil {
		return err
	}

	reaper := reaper.NewReaper(ctx)

	hook, err := otel.InitializeOpenTelemetry(ctx, cfg.OpenTelemetry)
	if err != nil {
		return err
	}
	reaper.AddContextErrorHook(hook)
	otel.NewTracer(cfg.OpenTelemetry)

	poiStore, err := store.NewStore(*cfg)
	if err != nil {
		return err
	}

	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Get("/dino/:name/actions", handlers.AllActionsHandler(*cfg, poiStore))
	app.Get("/dino/:name/actions/next", handlers.NextActionHandler(*cfg, poiStore))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			reaper.ErrCh <- err
		}
	}()

	reaper.AddErrorHook(app.Shutdown)

	reaper.Start(ctx)

	reaper.Wait()

	return nil
}
