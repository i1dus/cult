package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type AuthService interface {
	Login(ctx context.Context, phoneNumber string, password string) (uuid.UUID, string, error)
	RegisterNewUser(ctx context.Context, phoneNumber string, password string) (userID uuid.UUID, err error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, update domain.UserUpdate) error
}

type ParkingLotService interface {
	GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error)
	GetParkingLotByNumber(ctx context.Context, number string) (domain.ParkingLot, error)
	GetParkingLotsByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error)
	UpdateParkingLot(ctx context.Context, parkingLot string, update domain.ParkingLotUpdate) error
}

type BookingService interface {
	GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error)
	GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error)
	AddBooking(ctx context.Context, booking domain.Booking) error
}

type serverAPI struct {
	sso.UnimplementedParkingAPIServer

	auth       AuthService
	parkingLot ParkingLotService
	booking    BookingService
}

func Register(gRPCServer *grpc.Server, auth AuthService, parkingLot ParkingLotService, booking BookingService) {
	sso.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth, parkingLot: parkingLot, booking: booking})
}
