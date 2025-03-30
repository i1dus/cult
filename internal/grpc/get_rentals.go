package grpc

import (
	"context"
	"cult/internal/domain"
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

func (s *serverAPI) GetRental(ctx context.Context, in *desc.GetRentalRequest) (*desc.GetRentalResponse, error) {
	rental, err := s.rental.GetRental(ctx, in.ParkingLot)
	if err != nil {
		return nil, err
	}

	return &desc.GetRentalResponse{
		Rental: rentalsToApi([]domain.Rental{rental})[0],
	}, nil
}
