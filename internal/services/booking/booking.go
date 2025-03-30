package booking

import (
	"context"
	"cult/internal/domain"
	"fmt"
	"log/slog"
	"time"

	"cult/internal/lib/logger/sl"
)

type BookingService struct {
	log            *slog.Logger
	bookingRepo    Repository
	bookingTimeout time.Duration
}

type Repository interface {
	AddBooking(ctx context.Context, booking domain.Booking) error
	GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error)
	GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error)
	GetParkingLotsByFilter(ctx context.Context, filter domain.Filter) ([]domain.ParkingLot, error)
}

func NewBookingService(log *slog.Logger, repo Repository) *BookingService {
	return &BookingService{
		log:         log,
		bookingRepo: repo,
	}
}

func (s *BookingService) GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error) {
	const op = "BookingService.GetBookingsByFilter"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching all bookings")

	bookings, err := s.bookingRepo.GetBookingsByFilter(ctx, filter)
	if err != nil {
		log.Error("failed to get bookings", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved bookings", slog.Int("count", len(bookings)))
	return bookings, nil
}

func (s *BookingService) GetParkingLotsByFilter(ctx context.Context, filter domain.Filter) ([]domain.ParkingLot, error) {
	const op = "BookingService.GetParkingLotsByFilter"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching all parking lots")

	lots, err := s.bookingRepo.GetParkingLotsByFilter(ctx, filter)
	if err != nil {
		log.Error("failed to get bookings", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved bookings", slog.Int("count", len(lots)))
	return lots, nil
}

func (s *BookingService) GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error) {
	const op = "BookingService.GetBooking"

	log := s.log.With(slog.String("op", op))

	booking, err := s.bookingRepo.GetBooking(ctx, parkingLot)
	if err != nil {
		log.Error("failed to get booking", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved booking")
	return booking, nil
}

func (s *BookingService) AddBooking(ctx context.Context, booking domain.Booking) error {
	const op = "BookingService.GetBookingsByFilter"

	log := s.log.With(slog.String("op", op))

	log.Info("adding a booking")

	err := s.bookingRepo.AddBooking(ctx, booking)
	if err != nil {
		log.Error("failed to add booking", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
