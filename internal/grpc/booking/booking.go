package booking

import desc "cult/internal/gen/booking"

type Booking interface {
}

type parkingAPI struct {
	desc.UnimplementedParkingAPIServer
	booking Booking
}
