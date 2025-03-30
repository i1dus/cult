package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type AuthService interface {
	Login(ctx context.Context, phoneNumber string, password string) (uuid.UUID, string, error)
	RegisterNewUser(ctx context.Context, phoneNumber string, password string) (userID uuid.UUID, err error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, update domain.UserUpdate) error
	UserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
}

type ParkingLotService interface {
	GetAllParkingLots(ctx context.Context, userID uuid.UUID) ([]domain.ParkingLot, error)
	GetParkingLotByNumber(ctx context.Context, userID uuid.UUID, number string) (domain.ParkingLot, error)
	GetParkingLotsByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error)
	UpdateParkingLot(ctx context.Context, parkingLot string, update domain.ParkingLotUpdate) error
}

type BookingService interface {
	GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error)
	GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error)
	AddBooking(ctx context.Context, booking domain.Booking) error
	GetParkingLotsByFilter(ctx context.Context, filter domain.Filter) ([]domain.ParkingLot, error)
	EditBooking(ctx context.Context, bookingID uuid.UUID, to time.Time) error
	UpdateBookingVehicle(ctx context.Context, bookingID uuid.UUID, vehicle string) error
}

type RentalService interface {
	GetRentalsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Rental, error)
	AddRental(ctx context.Context, rental domain.Rental) error
}

type PaymentService interface {
}

type serverAPI struct {
	sso.UnimplementedParkingAPIServer

	auth       AuthService
	parkingLot ParkingLotService
	booking    BookingService
	rental     RentalService
	payment    PaymentService
}

func Register(gRPCServer *grpc.Server, auth AuthService, parkingLot ParkingLotService, booking BookingService, rental RentalService, payment PaymentService) {
	sso.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth, parkingLot: parkingLot, booking: booking, rental: rental, payment: payment})
}
