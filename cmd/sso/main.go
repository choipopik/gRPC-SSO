package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/choipopik/gRPC-SSO/internal/app"
	"github.com/choipopik/gRPC-SSO/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//config
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("App started...")

	//init app
	ssoApp := app.New(log, cfg.GRPC.Port, cfg.Storage_path, cfg.TokenTTL)

	go ssoApp.GRPCSrv.MustRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ssoApp.GRPCSrv.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

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
