package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
import desc "cult/pkg"

func (s *serverAPI) EditParkingBooking(ctx context.Context, in *desc.EditParkingBookingRequest) (*desc.EditParkingBookingResponse, error) {
	id, err := uuid.Parse(in.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.booking.EditBooking(ctx, id, in.GetTimeTo().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.EditParkingBookingResponse{}, nil
}

func (s *serverAPI) UpdateBookingVehicle(ctx context.Context, in *desc.UpdateBookingVehicleRequest) (*desc.UpdateBookingVehicleResponse, error) {
	id, err := uuid.Parse(in.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.booking.UpdateBookingVehicle(ctx, id, in.Vehicle)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.UpdateBookingVehicleResponse{}, nil
}
