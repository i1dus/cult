package grpc

import (
	"context"
	desc "cult/pkg"
)

func (s *serverAPI) GetRentals(ctx context.Context, in *desc.GetRentalsRequest) (*desc.GetRentalsResponse, error) {

	rentals, err := s.rental.GetRentalsByFilter(ctx, apiToFilter(in.GetFilter()))
	if err != nil {
		return nil, err
	}

	return &desc.GetRentalsResponse{
		Rentals: rentalsToApi(rentals),
		Total:   int64(len(rentals)),
	}, nil
}
