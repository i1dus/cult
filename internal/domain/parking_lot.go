package domain

import (
	sso "cult/pkg"
	"github.com/google/uuid"
)

type ParkingLot struct {
	ID          string
	ParkingType ParkingType
	VehicleID   *uuid.UUID
	OwnerID     *uuid.UUID
}

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

func (p ParkingType) GetPBType() sso.ParkingType {
	switch p {
	case UndefinedParkingType:
		return sso.ParkingType_UNDEFINED_PARKING_TYPE
	case PermanentParkingType:
		return sso.ParkingType_PERMANENT_PARKING_TYPE
	case RentParkingType:
		return sso.ParkingType_RENT_PARKING_TYPE
	case SpecialParkingType:
		return sso.ParkingType_SPECIAL_PARKING_TYPE
	case InclusiveParkingType:
		return sso.ParkingType_INCLUSIVE_PARKING_TYPE
	}
	panic("unknown parking type")
}
