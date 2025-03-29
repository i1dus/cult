package booking

import (
	"context"
	desc "cult/internal/gen/booking"
)

func (s *parkingAPI) GetBooking(ctx context.Context, in desc.AddParkingBookingRequest) (*desc.AddParkingBookingResponse, error) {
	return &desc.AddParkingBookingResponse{}, nil
}
