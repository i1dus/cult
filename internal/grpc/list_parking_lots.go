package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"
	"fmt"
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
				Number: int64(lo.Must(strconv.Atoi(item.ID))),
				Type:   item.ParkingType.GetPBType(),
				Status: getRandomParkingStatus(index),
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
