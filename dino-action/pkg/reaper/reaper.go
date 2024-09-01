package reaper

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Reaper struct {
	hooks     []func(context.Context, os.Signal)
	waitGroup sync.WaitGroup
	ErrCh     chan error
	logger    *zerolog.Logger
}

func NewReaper(ctx context.Context) *Reaper {
	logger := log.Ctx(ctx).With().Str("component", "reaper").Logger()

	return &Reaper{
		logger:    &logger,
		hooks:     []func(context.Context, os.Signal){},
		waitGroup: sync.WaitGroup{},
		ErrCh:     make(chan error),
	}
}

func (r *Reaper) AddContextSignalHook(hook func(context.Context, os.Signal)) {
	r.hooks = append(r.hooks, hook)
}

func (r *Reaper) AddContextErrorHook(hook func(context.Context) error) {
	r.hooks = append(r.hooks, func(ctx context.Context, _ os.Signal) {
		if err := hook(ctx); err != nil {
			r.logger.Error().Err(err).Msg("Hook failed")
		}
	})
}

func (r *Reaper) AddErrorHook(hook func() error) {
	r.hooks = append(r.hooks, func(ctx context.Context, _ os.Signal) {
		if err := hook(); err != nil {
			r.logger.Error().Err(err).Msg("Hook failed")
		}
	})
}

func (r *Reaper) AddContextHook(hook func(context.Context)) {
	r.hooks = append(r.hooks, func(ctx context.Context, _ os.Signal) {
		hook(ctx)
	})
}

func (r *Reaper) Start(ctx context.Context) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	r.waitGroup.Add(1)

	go func() {
		defer r.waitGroup.Done()
		select {
		case <-ctx.Done():
			r.runHooks(syscall.SIGSTOP)
		case err := <-r.ErrCh:
			r.logger.Error().Err(err).Msg("Error received")
			r.runHooks(syscall.SIGABRT)
		case s := <-sigCh:
			r.runHooks(s)
		}

		r.logger.Info().Msg("Shutdown complete")
	}()
}

func (r *Reaper) runHooks(s os.Signal) {
	r.logger.Warn().Str("signal", s.String()).Msg("Received signal, Shutting down all systems")
	ctx := context.Background()
	for _, hook := range r.hooks {
		hook(ctx, s)
	}
}

func (r *Reaper) Wait() {
	r.waitGroup.Wait()
}
