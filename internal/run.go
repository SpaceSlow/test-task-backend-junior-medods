package internal

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/application"
)

func RunServer() error {
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(rootCtx)
	context.AfterFunc(ctx, func() {
		timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelCtx()

		<-timeoutCtx.Done()
		slog.Error("failed to gracefully shutdown the service")
	})

	cfg := LoadServerConfig()

	httpServer := application.SetupHTTPServer()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: httpServer,
	}
	g.Go(func() error {
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), cfg.MaxTimeoutShutdown)
		defer cancelShutdownTimeoutCtx()
		return srv.Shutdown(shutdownTimeoutCtx)
	})

	return g.Wait()
}
