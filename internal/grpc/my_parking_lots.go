package grpc

import (
	"context"
	"cult/internal/domain"
	desc "cult/pkg"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) GetMyParkingLots(ctx context.Context, in *desc.GetMyParkingLotsRequest) (*desc.GetMyParkingLotsResponse, error) {
	parsedOwner := uuid.MustParse(in.GetFilter().OwnerId)

	rentals, err := s.rental.GetRentalsByFilter(ctx, apiToFilter(in.GetFilter()))
	if err != nil {
		return nil, err
	}

	bookings, err := s.booking.GetBookingsByFilter(ctx, apiToFilter(in.GetFilter()))
	if err != nil {
		return nil, err
	}

	parkingLots, err := s.parkingLot.GetParkingLotsByOwner(ctx, parsedOwner)
	if err != nil {
		return nil, err
	}

	var parkingLotIDs []int64
	for _, rental := range rentals {
		parkingLotIDs = append(parkingLotIDs, rental.ParkingLot)
	}
	for _, booking := range bookings {
		parkingLotIDs = append(parkingLotIDs, booking.ParkingLot)
	}

	resultLots := make([]domain.ParkingLot, 0)
	for _, lot := range parkingLots {
		if !lo.Contains(parkingLotIDs, lot.ID) {
			resultLots = append(resultLots, lot)
		}
	}

	return &desc.GetMyParkingLotsResponse{
		Rentals:     rentalsToApi(rentals),
		Bookings:    bookingsToApi(bookings),
		ParkingLots: parkingLotsToApi(resultLots),
		Total:       int64(len(rentals) + len(bookings) + len(parkingLots)),
	}, nil
}

func rentalsToApi(rentals []domain.Rental) []*desc.Rental {
	out := make([]*desc.Rental, 0, len(rentals))

	for _, rental := range rentals {
		out = append(
			out,
			&desc.Rental{
				ParkingLot:  rental.ParkingLot,
				TimeFrom:    timestamppb.New(rental.From),
				TimeTo:      timestamppb.New(rental.To),
				CostPerHour: rental.CostPerHour,
				CostPerDay:  rental.CostPerDay,
			},
		)
	}

	return out
}

func parkingLotsToApi(lots []domain.ParkingLot) []*desc.ParkingLot {
	out := make([]*desc.ParkingLot, 0, len(lots))

	for _, lot := range lots {
		out = append(
			out,
			&desc.ParkingLot{
				Number:  lot.ID,
				Vehicle: lot.OwnerVehicle,
			},
		)
	}

	return out
}
