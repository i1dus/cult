package grpc

import (
	"context"
	"cult/internal/domain"

	"github.com/google/uuid"
)
import desc "cult/pkg"

func (s *serverAPI) AddBooking(ctx context.Context, in *desc.AddParkingBookingRequest) (*desc.AddParkingBookingResponse, error) {

	err := s.booking.AddBooking(ctx, apiToBooking(in.GetBooking()))
	if err != nil {
		return nil, err
	}

	return &desc.AddParkingBookingResponse{
		ParkingLot: in.GetBooking().GetParkingLot(),
	}, nil
}

func apiToBooking(in *desc.ParkingBooking) domain.Booking {
	return domain.Booking{
		UserID:     uuid.MustParse(in.GetUserId()),
		ParkingLot: in.ParkingLot,
		From:       in.TimeFrom.AsTime(),
		To:         in.TimeTo.AsTime(),
		Vehicle:    in.Vehicle,
	}
}
