package grpc

import (
	"context"
	desc "cult/pkg"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) GetBookingPrice(ctx context.Context, in *desc.GetBookingPriceRequest) (*desc.GetBookingPriceResponse, error) {

	id, err := uuid.Parse(in.BookingId)

	price, err := s.rental.GetBookingPriceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if price == 0 {
		return nil, status.Error(codes.NotFound, "booking not found")
	}

	return &desc.GetBookingPriceResponse{
		AmountCents: price,
	}, nil
}
