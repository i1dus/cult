package parking_lot

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository/parking_lot"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"log/slog"
	"strconv"
	"time"

	"cult/internal/lib/logger/sl"
)

type ParkingLotService struct {
	log            *slog.Logger
	parkingLotRepo ParkingLotRepository
	bookingRepo    BookingRepository
	bookingTimeout time.Duration
}

type ParkingLotRepository interface {
	GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error)
	GetParkingLotByNumber(ctx context.Context, number string) (domain.ParkingLot, error)
	GetParkingLotsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error)
	UpdateParkingLot(ctx context.Context, parkingLotID string, update domain.ParkingLotUpdate) error
}

type BookingRepository interface {
	GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error)
}

func NewParkingLotService(log *slog.Logger, repo ParkingLotRepository, bookingRepo BookingRepository) *ParkingLotService {
	return &ParkingLotService{
		log:            log,
		bookingRepo:    bookingRepo,
		parkingLotRepo: repo,
	}
}

func (s *ParkingLotService) GetAllParkingLots(ctx context.Context, userID uuid.UUID) ([]domain.ParkingLot, error) {
	const op = "ParkingLotService.GetAllParkingLots"

	log := s.log.With(slog.String("op", op))

	log.Info("fetching all parking lots")

	lots, err := s.parkingLotRepo.GetAllParkingLots(ctx)
	if err != nil {
		log.Error("failed to get parking lots", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, lot := range lots {
		pType, pStatus := s.calculateTypeAndStatus(ctx, lot, userID)
		lot.ParkingType = pType
		lot.ParkingStatus = pStatus
	}

	log.Info("successfully retrieved parking lots", slog.Int("count", len(lots)))
	return lots, nil
}

func (s *ParkingLotService) GetParkingLotByNumber(ctx context.Context, userID uuid.UUID, number string) (domain.ParkingLot, error) {
	const op = "ParkingLotService.GetParkingLotByNumber"

	log := s.log.With(
		slog.String("op", op),
		slog.String("parking_number", number),
	)

	log.Info("fetching parking lot by number")

	lot, err := s.parkingLotRepo.GetParkingLotByNumber(ctx, number)
	if err != nil {
		if errors.Is(err, parking_lot.ErrNotFound) {
			log.Warn("parking lot not found", sl.Err(err))
			return domain.ParkingLot{}, fmt.Errorf("%s: %s", op, err.Error())
		}
		log.Error("failed to get parking lot", sl.Err(err))
		return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, err)
	}

	pType, pStatus := s.calculateTypeAndStatus(ctx, lot, userID)
	lot.ParkingType = pType
	lot.ParkingStatus = pStatus

	log.Info("successfully retrieved parking lot")
	return lot, nil
}

func (s *ParkingLotService) GetParkingLotsByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error) {
	const op = "ParkingLotService.GetParkingLotsByOwner"

	log := s.log.With(
		slog.String("op", op),
		slog.String("owner_id", ownerID.String()),
	)

	log.Info("fetching parking lots by owner")

	lots, err := s.parkingLotRepo.GetParkingLotsByOwnerID(ctx, ownerID)
	if err != nil {
		log.Error("failed to get parking lots by owner", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, lot := range lots {
		pType, pStatus := s.calculateTypeAndStatus(ctx, lot, ownerID)
		lot.ParkingType = pType
		lot.ParkingStatus = pStatus
	}

	log.Info("successfully retrieved parking lots by owner",
		slog.Int("count", len(lots)),
	)
	return lots, nil
}

func (s *ParkingLotService) UpdateParkingLot(ctx context.Context, parkingLot string, update domain.ParkingLotUpdate) error {
	return s.parkingLotRepo.UpdateParkingLot(ctx, parkingLot, update)
}

func (s *ParkingLotService) calculateTypeAndStatus(ctx context.Context, lot domain.ParkingLot, userID uuid.UUID) (domain.ParkingType, domain.ParkingLotStatus) {
	if lot.ParkingKind != domain.RegularParkingKind {
		return domain.UndefinedParkingType, domain.UndefinedParkingLotStatus
	}

	// Имеет ли парковочное место владельца
	var hasOwner bool
	// Является ли юзер владельцем
	var isOwner bool

	if lot.OwnerID == nil {
		hasOwner = true
		if *lot.OwnerID == userID {
			isOwner = true
		}
	}

	idInt := int64(lo.Must(strconv.Atoi(lot.ID)))
	booking, err := s.bookingRepo.GetBooking(ctx, idInt)
	if err != nil {
		s.log.Error(err.Error())
		return domain.UndefinedParkingType, domain.UndefinedParkingLotStatus
	}

	if booking == nil {
		// Брони нет
		if isOwner {
			return domain.UndefinedParkingType, domain.MineParkingLotStatus
		}
		if hasOwner {
			return domain.UndefinedParkingType, domain.BusyParkingLotStatus
		}
		return domain.UndefinedParkingType, domain.FreeParkingLotStatus
	}

	if booking.UserID != userID {
		// Бронь не наша
		if !isOwner {
			return domain.OwnedParkingType, domain.BusyParkingLotStatus
		}

		if booking.IsShortTerm {
			return domain.ShortTermRentByOtherParkingType, domain.MineParkingLotStatus
		}
		return domain.LongTermRentByOtherParkingType, domain.MineParkingLotStatus
	}

	// Бронь наша
	status := domain.MineParkingLotStatus
	if time.Now().After(booking.To) {
		status = domain.ExpiredParkingLotStatus
	}

	if booking.IsShortTerm {
		return domain.ShortTermRentByMeParkingType, status
	}
	return domain.LongTermRentByMeParkingType, status
}
