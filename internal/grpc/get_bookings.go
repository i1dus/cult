package grpc

import (
	"context"
	"cult/internal/domain"
	desc "cult/pkg"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) GetParkingBookingsList(ctx context.Context, in *desc.GetParkingBookingsListRequest) (*desc.GetParkingBookingsListResponse, error) {

	bookings, err := s.booking.GetBookingsByFilter(ctx, apiToFilter(in.Filter))
	if err != nil {
		return nil, err
	}
	return &desc.GetParkingBookingsListResponse{
		Bookings: bookingsToApi(bookings),
		Total:    int64(len(bookings)),
	}, nil
}

func bookingsToApi(bookings []domain.Booking) []*desc.ParkingBooking {
	out := make([]*desc.ParkingBooking, 0, len(bookings))

	for _, booking := range bookings {
		out = append(out, &desc.ParkingBooking{
			UserId:     booking.UserID.String(),
			ParkingLot: booking.ParkingLot,
			Vehicle:    booking.Vehicle,
			TimeFrom:   timestamppb.New(booking.From),
			TimeTo:     timestamppb.New(booking.To),
			//ParkingLot: booking.ParkingLot,
			BookingId: booking.ID.String(),
		})
	}

	return out
}

func apiToFilter(in *desc.Filter) domain.Filter {
	return domain.Filter{
		UserID:      uuid.MustParse(in.GetOwnerId()),
		From:        in.TimeFrom.AsTime(),
		To:          in.TimeTo.AsTime(),
		ParkingLots: nil,
	}
}

func bookingToApi(booking *domain.Booking) *desc.ParkingBooking {
	return &desc.ParkingBooking{
		UserId:     booking.UserID.String(),
		ParkingLot: booking.ParkingLot,
		Vehicle:    booking.Vehicle,
		TimeFrom:   timestamppb.New(booking.From),
		TimeTo:     timestamppb.New(booking.To),
		//ParkingLot: booking.ParkingLot,
		BookingId: booking.ID.String(),
	}
}
