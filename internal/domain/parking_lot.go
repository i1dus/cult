package domain

type ParkingLotStatus string

const (
	UndefinedParkingLotStatus  ParkingLotStatus = "UNDEFINED_PARKING_LOT_STATUS"
	AvailableParkingLotStatus  ParkingLotStatus = "AVAILABLE_PARKING_LOT_STATUS"
	BookedParkingLotStatus     ParkingLotStatus = "BOOKED_PARKING_LOT_STATUS"
	BookedByMeParkingLotStatus ParkingLotStatus = "BOOKED_BY_ME_PARKING_LOT_STATUS"
)

type ParkingType string

const (
	UndefinedParkingType ParkingType = "UNDEFINED_PARKING_TYPE"
	PermanentParkingType ParkingType = "PERMANENT_PARKING_TYPE"
	RentParkingType      ParkingType = "RENT_PARKING_TYPE"
	SpecialParkingType   ParkingType = "SPECIAL_PARKING_TYPE"
	InclusiveParkingType ParkingType = "INCLUSIVE_PARKING_TYPE"
)
