package app

import (
	grpcapp "cult/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	//storage, err := sqlite.New(storagePath)
	//if err != nil {
	//	panic(err)
	//}
	//
	//authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, nil, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
