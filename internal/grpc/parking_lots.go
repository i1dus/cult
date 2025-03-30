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
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is empty")
	}

	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse user ID: %s", err.Error())
	}

	lots, err := s.parkingLot.GetAllParkingLots(ctx, userID)
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
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is empty")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse user ID: %s", err.Error())
	}

	parkingNumber := strconv.FormatInt(req.Number, 10)

	lot, err := s.parkingLot.GetParkingLotByNumber(ctx, userID, parkingNumber)
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
	var vehicle *string
	if item.CurrentVehicle != nil {
		vehicle = item.CurrentVehicle
	} else if item.OwnerVehicle != nil {
		vehicle = item.OwnerVehicle
	}

	var ownerID *string
	if item.OwnerID != nil {
		ownerID = lo.ToPtr(item.OwnerID.String())
	}

	return &sso.ParkingLot{
		Number: item.ID,
		Kind:    item.ParkingKind.GetPBType(),
		Type:    item.ParkingType.GetPBType(),
		Status:  item.ParkingStatus.ParkingLotStatusToPB(),
		OwnerId: ownerID,
		Vehicle: vehicle,
	}
}
