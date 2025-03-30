package grpc

import (
	"context"
	"cult/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import desc "cult/pkg"

func (s *serverAPI) AddRental(ctx context.Context, in *desc.AddRentalRequest) (*desc.AddRentalResponse, error) {
	err := s.rental.AddRental(ctx, apiToRental(in.GetRental()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.AddRentalResponse{
		ParkingLotId: in.GetRental().GetParkingLot(),
	}, nil
}

func apiToRental(in *desc.Rental) domain.Rental {
	return domain.Rental{
		ParkingLot:  in.ParkingLot,
		From:        in.TimeFrom.AsTime(),
		To:          in.TimeTo.AsTime(),
		CostPerHour: in.CostPerHour,
		CostPerDay:  in.CostPerDay,
	}
}
