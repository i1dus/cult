package grpc

import (
	"context"
	desc "cult/pkg"
)

func (s *parkingAPI) GetBooking(ctx context.Context, in *desc.AddParkingBookingRequest) (*desc.AddParkingBookingResponse, error) {
	return &desc.AddParkingBookingResponse{}, nil
}
