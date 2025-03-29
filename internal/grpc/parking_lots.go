package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (s *serverAPI) ListParkingLots(ctx context.Context, in *sso.ListParkingLotsRequest) (*sso.ListParkingLotsResponse, error) {
	lots, err := s.parkingLot.GetAllParkingLots(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get list parking lots: %s", err.Error()))
	}

	return &sso.ListParkingLotsResponse{
		ParkingLot: lo.Map(lots, func(item domain.ParkingLot, index int) *sso.ParkingLot {
			return &sso.ParkingLot{
				Number:  int64(lo.Must(strconv.Atoi(item.ID))),
				Type:    item.ParkingType.GetPBType(),
				Status:  getRandomParkingStatus(index),
				OwnerId: lo.ToPtr(item.OwnerID.String()),
			}
		}),
		Total: int64(len(lots)),
	}, nil
}

func getRandomParkingStatus(index int) sso.ParkingLotStatus {
	switch index % 3 {
	case 0:
		return sso.ParkingLotStatus_AVAILABLE_PARKING_LOT_STATUS
	case 1:
		return sso.ParkingLotStatus_BOOKED_PARKING_LOT_STATUS
	case 2:
		return sso.ParkingLotStatus_BOOKED_BY_ME_PARKING_LOT_STATUS
	}
	panic("undefined parking lot status")
}

func (s *serverAPI) GetParkingLot(ctx context.Context, req *sso.GetParkingLotRequest) (*sso.GetParkingLotResponse, error) {
	parkingNumber := strconv.FormatInt(req.Number, 10)

	lot, err := s.parkingLot.GetParkingLotByNumber(ctx, parkingNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get parking lot: %v", err))
	}

	return &sso.GetParkingLotResponse{
		ParkingLot: &sso.ParkingLot{
			Number:  req.Number,
			Type:    lot.ParkingType.GetPBType(),
			Status:  getRandomParkingStatus(int(req.Number) % 3),
			OwnerId: lo.ToPtr(lot.OwnerID.String()),
		},
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
			number := lo.Must(strconv.ParseInt(item.ID, 10, 64))

			return &sso.ParkingLot{
				Number:  number,
				Type:    item.ParkingType.GetPBType(),
				Status:  getRandomParkingStatus(index),
				OwnerId: lo.ToPtr(item.OwnerID.String()),
			}
		}),
	}, nil
}
