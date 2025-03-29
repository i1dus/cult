package grpc

import desc "cult/pkg"

type Booking interface {
}

type parkingAPI struct {
	desc.UnimplementedParkingAPIServer
	booking Booking
}
