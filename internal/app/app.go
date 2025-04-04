package app

import (
	grpcapp "cult/internal/app/grpc"
	"cult/internal/repository/booking_repository"
	parking_lot2 "cult/internal/repository/parking_lot"
	payment_repository "cult/internal/repository/payment"
	"cult/internal/repository/rental_repository"
	"cult/internal/repository/user_repository"
	"cult/internal/services/auth"
	"cult/internal/services/booking"
	"cult/internal/services/parking_lot"
	"cult/internal/services/payment"
	"cult/internal/services/rental"
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
	rentalRepository := rental_repository.NewRentalRepository(conn, log)
	paymentRepository := payment_repository.NewRepository()
	paymentProcessor := payment.NewSberProcessor(
		"https://3dsec.sberbank.ru",
		"admin",
		"password",
	)

	authService := auth.New(log, userRepository, tokenTTL, secret)
	parkingLotService := parking_lot.NewParkingLotService(log, parkingLotRepository, bookingRepository, userRepository)
	bookingService := booking.NewBookingService(log, bookingRepository)
	rentalService := rental.NewRentalService(log, rentalRepository)
	paymentService := payment.NewService(paymentRepository, paymentProcessor)

	grpcApp := grpcapp.New(log, authService, parkingLotService, bookingService, rentalService, paymentService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
