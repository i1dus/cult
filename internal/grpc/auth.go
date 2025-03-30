package grpc

import (
	"context"
	"cult/internal/domain"
	desc "cult/pkg"

	"github.com/google/uuid"

	"google.golang.org/grpc"
)

type AuthService interface {
	Login(ctx context.Context, phoneNumber string, password string) (uuid.UUID, string, error)
	RegisterNewUser(ctx context.Context, phoneNumber string, password string) (userID uuid.UUID, err error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type ParkingLotService interface {
	GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error)
}

type BookingService interface {
	GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error)
	GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error)
	AddBooking(ctx context.Context, booking domain.Booking) error
}

type serverAPI struct {
	desc.UnimplementedParkingAPIServer

	auth       AuthService
	parkingLot ParkingLotService
	booking    BookingService
}

func Register(gRPCServer *grpc.Server, auth AuthService, parkingLot ParkingLotService, booking BookingService) {
	desc.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth, parkingLot: parkingLot, booking: booking})
}
