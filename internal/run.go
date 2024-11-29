package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/application"
	usersrepo "github.com/SpaceSlow/test-task-backend-junior-medods/internal/infrastructure/users"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/service/notifier"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/service/users"
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

	repo, err := usersrepo.NewPostgresRepo(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to initialize a user repo: %w", err)
	}
	defer repo.Close()

	g.Go(func() error {
		defer slog.Info("closed user repo")
		<-ctx.Done()
		repo.Close()
		return nil
	})

	notifierService := notifier.NewSMTPNotifierService(cfg)

	userService := users.NewUserService(repo, notifierService, cfg)

	httpServer := application.SetupHTTPServer(userService)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: httpServer,
	}
	g.Go(func() error {
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		defer slog.Info("stopped http server")
		<-ctx.Done()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), cfg.MaxTimeoutShutdown)
		defer cancelShutdownTimeoutCtx()
		return srv.Shutdown(shutdownTimeoutCtx)
	})

	return g.Wait()
}
