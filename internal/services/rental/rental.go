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
}

func NewRentalService(log *slog.Logger, repo Repository) *RentalService {
	return &RentalService{
		log:        log,
		rentalRepo: repo,
	}
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
