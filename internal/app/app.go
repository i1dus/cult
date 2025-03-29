package app

import (
	grpcapp "cult/internal/app/grpc"
	"cult/internal/repository/user_repository"
	"cult/internal/services/auth"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger, grpcPort int,
	conn *pgx.Conn, tokenTTL time.Duration, secret string) *App {

	userRepository := user_repository.NewUserRepository(conn, log)

	authService := auth.New(log, userRepository, tokenTTL, secret)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
