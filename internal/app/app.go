package app

import (
	grpcapp "cult/internal/app/grpc"
	"cult/internal/repository/booking_repository"
	parking_lot2 "cult/internal/repository/parking_lot"
	"cult/internal/repository/user_repository"
	"cult/internal/services/auth"
	"cult/internal/services/booking"
	"cult/internal/services/parking_lot"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger, grpcPort int,
	conn *pgx.Conn, tokenTTL time.Duration, secret string) *App {

	userRepository := user_repository.NewUserRepository(conn, log)
	parkingLotRepository := parking_lot2.NewParkingLotRepository(conn, log)

	bookingRepository := booking_repository.NewBookingRepository(conn, log)

	authService := auth.New(log, userRepository, tokenTTL, secret)
	parkingLotService := parking_lot.NewParkingLotService(log, parkingLotRepository)
	bookingService := booking.NewBookingService(log, bookingRepository)

	grpcApp := grpcapp.New(log, authService, parkingLotService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
