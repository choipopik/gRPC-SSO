package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/choipopik/gRPC-SSO/internal/app/grpc"
	"github.com/choipopik/gRPC-SSO/internal/services/auth"
	"github.com/choipopik/gRPC-SSO/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	//init db
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	//init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, port, authService)

	return &App{
		grpcApp,
	}
}
