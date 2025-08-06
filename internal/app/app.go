package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/choipopik/gRPC-SSO/internal/app/grpc"
	"github.com/choipopik/gRPC-SSO/internal/services/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	//init db

	//init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, port, authService)

	return &App{
		grpcApp,
	}
}
