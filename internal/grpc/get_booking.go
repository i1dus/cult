package grpc

import (
	"context"
	desc "cult/pkg"
)

func (s *serverAPI) GetBooking(ctx context.Context, in *desc.GetParkingBookingRequest) (*desc.GetParkingBookingResponse, error) {

	booking, err := s.booking.GetBooking(ctx, in.GetParkingLot())
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return &desc.GetParkingBookingResponse{}, nil
	}

	return &desc.GetParkingBookingResponse{
		Booking: bookingToApi(booking),
	}, nil
}
