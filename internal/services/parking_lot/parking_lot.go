package parking_lot

import (
	"context"
	"cult/internal/domain"
	"fmt"
	"log/slog"
	"time"

	"cult/internal/lib/logger/sl"
)

type ParkingLotService struct {
	log            *slog.Logger
	parkingLotRepo ParkingLotRepository
	bookingTimeout time.Duration
}

type ParkingLotRepository interface {
	GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error)
}

func NewParkingLotService(log *slog.Logger, repo ParkingLotRepository) *ParkingLotService {
	return &ParkingLotService{
		log:            log,
		parkingLotRepo: repo,
	}
}

func (s *ParkingLotService) GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error) {
	const op = "ParkingLotService.GetAllParkingLots"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching all parking lots")

	lots, err := s.parkingLotRepo.GetAllParkingLots(ctx)
	if err != nil {
		log.Error("failed to get parking lots", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully retrieved parking lots", slog.Int("count", len(lots)))
	return lots, nil
}

//// GetParkingLotByID returns a single parking lot by its ID
//func (s *ParkingLotService) GetParkingLotByID(ctx context.Context, id string) (domain.ParkingLot, error) {
//	const op = "ParkingLotService.GetParkingLotByID"
//
//	log := s.log.With(
//		slog.String("op", op),
//		slog.String("parking_lot_id", id),
//	)
//
//	log.Info("fetching parking lot")
//
//	lot, err := s.parkingLotRepo.GetByID(ctx, id)
//	if err != nil {
//		if errors.Is(err, domain.ErrNotFound) {
//			log.Warn("parking lot not found", sl.Err(err))
//			return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, ErrParkingLotNotFound)
//		}
//		log.Error("failed to get parking lot", sl.Err(err))
//		return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, err)
//	}
//
//	log.Info("successfully retrieved parking lot")
//	return lot, nil
//}
