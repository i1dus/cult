package app

import (
	"context"
	grpcapp "cult/internal/app/grpc"
	"cult/internal/services/auth"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	databaseURL string,
	tokenTTL time.Duration,
) *App {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	authService := auth.New(log, nil, nil, nil, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
