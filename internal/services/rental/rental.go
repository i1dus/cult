package rental

import (
	"context"
	"cult/internal/domain"
	"cult/internal/lib/logger/sl"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type RentalService struct {
	log        *slog.Logger
	rentalRepo Repository
}

type Repository interface {
	AddRental(ctx context.Context, rental domain.Rental) error
	GetBookingPriceByID(ctx context.Context, bookingID uuid.UUID) (int64, error)
	GetRentalsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Rental, error)
	GetRentalByLot(ctx context.Context, parkingLot int64) (domain.Rental, error)
}

func NewRentalService(log *slog.Logger, repo Repository) *RentalService {
	return &RentalService{
		log:        log,
		rentalRepo: repo,
	}
}

func (s *RentalService) GetBookingPriceByID(ctx context.Context, bookingID uuid.UUID) (int64, error) {
	const op = "RentalService.GetBookingPriceByID"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching price")
	price, err := s.rentalRepo.GetBookingPriceByID(ctx, bookingID)
	if err != nil {
		log.Error("failed to get price", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved rentals")
	return price, nil
}

func (s *RentalService) GetRentalsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Rental, error) {
	const op = "RentalService.GetRentalsByFilter"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching all rentals")

	rentals, err := s.rentalRepo.GetRentalsByFilter(ctx, filter)
	if err != nil {
		log.Error("failed to get rentals", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved rentals", slog.Int("count", len(rentals)))
	return rentals, nil
}

func (s *RentalService) AddRental(ctx context.Context, rental domain.Rental) error {
	const op = "RentalService.AddRental"

	log := s.log.With(slog.String("op", op))

	log.Info("adding a rental")

	err := s.rentalRepo.AddRental(ctx, rental)
	if err != nil {
		log.Error("failed to add rental", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RentalService) GetRental(ctx context.Context, parkingLot int64) (domain.Rental, error) {
	const op = "RentalService.GetRental"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching rental")

	rental, err := s.rentalRepo.GetRentalByLot(ctx, parkingLot)
	if err != nil {
		log.Error("failed to get rental", sl.Err(err))
		return domain.Rental{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved rental")
	return rental, nil
}
