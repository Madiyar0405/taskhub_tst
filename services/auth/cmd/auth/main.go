package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/lib/logger/handlers/slogpretty"
	"context"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	cfg = config.MustLoad()
	wg  sync.WaitGroup
)

func main() {
	log := setupLogger(cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(2)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.ServicePort),
	)

	application := app.New(log, cfg.ServicePort, cfg.HttpPort, cfg.PostgresURI, cfg.TokenTTL)

	go func() {
		defer wg.Done()
		application.GRPCSrv.MustRun()
	}()

	go func() {
		defer wg.Done()
		application.HTTPSrv.MustRun(ctx)
	}()

	go func() {
		c := cron.New()
		defer c.Stop()

		_, err := c.AddFunc("@every 24h", func() {
			err := application.Storage.DeleteInactiveUsers(context.Background())
			if err != nil {
				log.Error("error deleting inactive users", slog.String("error", err.Error()))
			} else {
				log.Info("unconfirmed users deleted successfully")
			}
		})
		if err != nil {
			log.Error("Failed to add cron job", slog.String("error", err.Error()))
		}

		c.Start()
		select {}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-quit

	log.Info("stopping application", slog.String("signal", s.String()))

	application.HTTPSrv.Stop(ctx)
	application.GRPCSrv.Stop()

	wg.Wait()

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
