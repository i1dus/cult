package grpc

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	sso "cult/pkg"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) ListParkingLots(ctx context.Context, in *sso.ListParkingLotsRequest) (*sso.ListParkingLotsResponse, error) {
	lots, err := s.parkingLot.GetAllParkingLots(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get list parking lots: %s", err.Error()))
	}

	return &sso.ListParkingLotsResponse{
		ParkingLot: lo.Map(lots, func(item domain.ParkingLot, index int) *sso.ParkingLot {
			return ConvertParkingLotPbToDomain(item, index)
		}),
		Total: int64(len(lots)),
	}, nil
}

func (s *serverAPI) GetParkingLot(ctx context.Context, req *sso.GetParkingLotRequest) (*sso.GetParkingLotResponse, error) {
	parkingNumber := strconv.FormatInt(req.Number, 10)

	lot, err := s.parkingLot.GetParkingLotByNumber(ctx, parkingNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get parking lot: %v", err))
	}

	return &sso.GetParkingLotResponse{
		ParkingLot: ConvertParkingLotPbToDomain(lot, 0),
	}, nil
}

func (s *serverAPI) GetParkingLotsByUserID(ctx context.Context, req *sso.GetParkingLotsByUserIDRequest) (*sso.GetParkingLotsByUserIDResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	lots, err := s.parkingLot.GetParkingLotsByOwner(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user parking lots: %v", err))
	}

	return &sso.GetParkingLotsByUserIDResponse{
		ParkingLot: lo.Map(lots, func(item domain.ParkingLot, index int) *sso.ParkingLot {
			return ConvertParkingLotPbToDomain(item, index)
		}),
	}, nil
}

func (s *serverAPI) UpdateParkingLot(ctx context.Context, req *sso.UpdateParkingLotRequest) (*sso.UpdateParkingLotResponse, error) {
	parkingLotID := strconv.FormatInt(req.Number, 10)

	updateLot := domain.ParkingLotUpdate{}

	if req.Kind != nil {
		pt := domain.ParkingKindFromPB(*req.Kind)
		if pt == domain.UndefinedParkingKind {
			return nil, status.Error(codes.InvalidArgument, "invalid parking kind")
		}
		updateLot.ParkingKind = &pt
	}

	if req.OwnerId != nil {
		ownerID, err := uuid.Parse(*req.OwnerId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid owner ID format")
		}
		updateLot.OwnerID = &ownerID
	}

	if req.OwnerVehicle != nil {
		updateLot.OwnerVehicle = req.OwnerVehicle
	}

	// Ничего не передано
	if updateLot.ParkingKind == nil && updateLot.OwnerID == nil && updateLot.OwnerVehicle == nil {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	err := s.parkingLot.UpdateParkingLot(ctx, parkingLotID, updateLot)
	if err != nil {
		if errors.Is(err, repository.ErrParkingLotNotFound) {
			return nil, status.Error(codes.NotFound, "parking lot not found")
		}
		if errors.Is(err, repository.ErrNoFieldsToUpdate) {
			return nil, status.Error(codes.InvalidArgument, "no fields to update")
		}
		return nil, status.Error(codes.Internal, "failed to update parking lot")
	}

	return &sso.UpdateParkingLotResponse{}, nil
}

func ConvertParkingLotPbToDomain(item domain.ParkingLot, index int) *sso.ParkingLot {
	return &sso.ParkingLot{
		Number: item.ID,
		Type:   item.ParkingType.GetPBType(),
		Status: getRandomParkingStatus(index),
		OwnerId: func(id *uuid.UUID) *string {
			if id == nil {
				return nil
			}
			return lo.ToPtr(id.String())
		}(item.OwnerID),
		Vehicle: item.OwnerVehicle,
	}
}

func getRandomParkingStatus(index int) sso.ParkingLotStatus {
	switch index % 4 {
	case 0:
		return sso.ParkingLotStatus_FREE_PARKING_LOT_STATUS
	case 1:
		return sso.ParkingLotStatus_BUSY_PARKING_LOT_STATUS
	case 2:
		return sso.ParkingLotStatus_MINE_PARKING_LOT_STATUS
	case 3:
		return sso.ParkingLotStatus_EXPIRED_PARKING_LOT_STATUS

	}

	panic("undefined parking lot status")
}
