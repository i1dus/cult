package grpc

import (
	"context"
	"cult/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import desc "cult/pkg"

func (s *serverAPI) AddParkingBooking(ctx context.Context, in *desc.AddParkingBookingRequest) (*desc.AddParkingBookingResponse, error) {

	fmt.Println(in.Booking.GetParkingLot())

	err := s.booking.AddBooking(ctx, apiToBooking(in.GetBooking()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.AddParkingBookingResponse{
		ParkingLot: in.GetBooking().GetParkingLot(),
	}, nil
}

func apiToBooking(in *desc.ParkingBooking) domain.Booking {
	fmt.Println("whhtth " + in.GetUserId())
	return domain.Booking{
		UserID:     uuid.MustParse(in.GetUserId()),
		ParkingLot: in.ParkingLot,
		From:       in.TimeFrom.AsTime(),
		To:         in.TimeTo.AsTime(),
		Vehicle:    in.Vehicle,
	}
}
