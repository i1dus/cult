package app

import (
	grpcapp "cult/internal/app/grpc"
	parking_lot2 "cult/internal/repository/parking_lot"
	"cult/internal/repository/user_repository"
	"cult/internal/services/auth"
	"cult/internal/services/parking_lot"
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
	parkingLotRepository := parking_lot2.NewParkingLotRepository(conn, log)

	authService := auth.New(log, userRepository, tokenTTL, secret)
	parkingLotService := parking_lot.NewParkingLotService(log, parkingLotRepository)

	grpcApp := grpcapp.New(log, authService, parkingLotService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
