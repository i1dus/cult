package domain

import (
	sso "cult/pkg"
	"time"

	"github.com/google/uuid"
)

type ParkingLot struct {
	ID             int64
	ParkingKind    ParkingKind
	ParkingType    ParkingType
	ParkingStatus  ParkingLotStatus
	OwnerID        *uuid.UUID
	OwnerVehicle   *string
	CurrentVehicle *string
}

type ParkingLotStatus string

const (
	UndefinedParkingLotStatus ParkingLotStatus = ""
	FreeParkingLotStatus      ParkingLotStatus = "FREE_PARKING_LOT_STATUS"
	BusyParkingLotStatus      ParkingLotStatus = "BUSY_PARKING_LOT_STATUS"
	// MineParkingLotStatus пользователь владеет ИЛИ арендует
	MineParkingLotStatus    ParkingLotStatus = "MINE_PARKING_LOT_STATUS"
	ExpiredParkingLotStatus ParkingLotStatus = "EXPIRED_PARKING_LOT_STATUS"
)

func (p ParkingLotStatus) ParkingLotStatusToPB() sso.ParkingLotStatus {
	switch p {
	case FreeParkingLotStatus:
		return sso.ParkingLotStatus_FREE_PARKING_LOT_STATUS
	case BusyParkingLotStatus:
		return sso.ParkingLotStatus_BUSY_PARKING_LOT_STATUS
	case MineParkingLotStatus:
		return sso.ParkingLotStatus_MINE_PARKING_LOT_STATUS
	case ExpiredParkingLotStatus:
		return sso.ParkingLotStatus_EXPIRED_PARKING_LOT_STATUS
	default:
		return sso.ParkingLotStatus_UNDEFINED_PARKING_LOT_STATUS
	}
}

func ParkingLotStatusFromPB(pbStatus sso.ParkingLotStatus) ParkingLotStatus {
	switch pbStatus {
	case sso.ParkingLotStatus_FREE_PARKING_LOT_STATUS:
		return FreeParkingLotStatus
	case sso.ParkingLotStatus_BUSY_PARKING_LOT_STATUS:
		return BusyParkingLotStatus
	case sso.ParkingLotStatus_MINE_PARKING_LOT_STATUS:
		return MineParkingLotStatus
	case sso.ParkingLotStatus_EXPIRED_PARKING_LOT_STATUS:
		return ExpiredParkingLotStatus
	default:
		return UndefinedParkingLotStatus
	}
}

type ParkingType string

const (
	// UndefinedParkingType неизвестный тип
	UndefinedParkingType ParkingType = ""
	// OwnedParkingType владеет кто-то другой
	OwnedParkingType ParkingType = "OWNED_PARKING_TYPE"
	// LongTermRentByMeParkingType чужое место, арендую долгосрочно
	LongTermRentByMeParkingType ParkingType = "LONG_TERM_RENT_BY_ME_PARKING_TYPE"
	// ShortTermRentByMeParkingType чужое место, арендую краткосрочно
	ShortTermRentByMeParkingType ParkingType = "SHORT_TERM_RENT_BY_ME_PARKING_TYPE"
	// LongTermRentByOtherParkingType мое место, сняли долгосрочно (в текущий момент)
	LongTermRentByOtherParkingType ParkingType = "LONG_TERM_RENT_BY_OTHER_PARKING_TYPE"
	// ShortTermRentByOtherParkingType мое место, сняли краткосрочно (в текущий момент)
	ShortTermRentByOtherParkingType ParkingType = "SHORT_TERM_RENT_BY_OTHER_PARKING_TYPE"
)

func (p ParkingType) String() string {
	return string(p)
}

// GetPBType - domain to proto
func (p ParkingType) GetPBType() sso.ParkingType {
	switch p {
	case OwnedParkingType:
		return sso.ParkingType_OWNED_PARKING_TYPE
	case LongTermRentByMeParkingType:
		return sso.ParkingType_LONG_TERM_RENT_BY_ME_PARKING_TYPE
	case ShortTermRentByMeParkingType:
		return sso.ParkingType_SHORT_TERM_RENT_BY_ME_PARKING_TYPE
	case LongTermRentByOtherParkingType:
		return sso.ParkingType_LONG_TERM_RENT_BY_OTHER_PARKING_TYPE
	case ShortTermRentByOtherParkingType:
		return sso.ParkingType_SHORT_TERM_RENT_BY_OTHER_PARKING_TYPE
	default:
		return sso.ParkingType_UNDEFINED_PARKING_TYPE
	}
}

// ParkingTypeFromPB - proto to domain
func ParkingTypeFromPB(pbType sso.ParkingType) ParkingType {
	switch pbType {
	case sso.ParkingType_OWNED_PARKING_TYPE:
		return OwnedParkingType
	case sso.ParkingType_LONG_TERM_RENT_BY_ME_PARKING_TYPE:
		return LongTermRentByMeParkingType
	case sso.ParkingType_SHORT_TERM_RENT_BY_ME_PARKING_TYPE:
		return ShortTermRentByMeParkingType
	case sso.ParkingType_LONG_TERM_RENT_BY_OTHER_PARKING_TYPE:
		return LongTermRentByOtherParkingType
	case sso.ParkingType_SHORT_TERM_RENT_BY_OTHER_PARKING_TYPE:
		return ShortTermRentByOtherParkingType
	default:
		return UndefinedParkingType
	}
}

type Booking struct {
	UserID      uuid.UUID
	RentalID    uuid.UUID
	From        time.Time
	To          time.Time
	Vehicle     string
	IsPresent   bool
	IsShortTerm bool
}

type Rental struct {
	ID          uuid.UUID
	ParkingLot  int64
	From        time.Time
	To          time.Time
	CostPerHour int64
	CostPerDay  int64
}

type Filter struct {
	UserID      uuid.UUID
	From        time.Time
	To          time.Time
	ParkingLots []int
}

type ParkingLotUpdate struct {
	ParkingKind  *ParkingKind
	OwnerID      *uuid.UUID
	OwnerVehicle *string
}

// ParkingKind вид места
type ParkingKind string

const (
	// UndefinedParkingKind неизвестный вид
	UndefinedParkingKind ParkingKind = ""
	// RegularParkingKind обычное место
	RegularParkingKind ParkingKind = "REGULAR_PARKING_KIND"
	// SpecialParkingKind место под спец. службы (скорая, пожарные...)
	SpecialParkingKind ParkingKind = "SPECIAL_PARKING_KIND"
	// InclusiveParkingKind место под инвалидов
	InclusiveParkingKind ParkingKind = "INCLUSIVE_PARKING_KIND"
)

func (p ParkingKind) String() string {
	return string(p)
}

// ParkingKindFromPB - proto to domain
func ParkingKindFromPB(pbKind sso.ParkingKind) ParkingKind {
	switch pbKind {
	case sso.ParkingKind_REGULAR_PARKING_KIND:
		return RegularParkingKind
	case sso.ParkingKind_SPECIAL_PARKING_KIND:
		return SpecialParkingKind
	case sso.ParkingKind_INCLUSIVE_PARKING_KIND:
		return InclusiveParkingKind
	default:
		return UndefinedParkingKind
	}
}

// GetPBType - domain to proto
func (p ParkingKind) GetPBType() sso.ParkingKind {
	switch p {
	case RegularParkingKind:
		return sso.ParkingKind_REGULAR_PARKING_KIND
	case SpecialParkingKind:
		return sso.ParkingKind_SPECIAL_PARKING_KIND
	case InclusiveParkingKind:
		return sso.ParkingKind_INCLUSIVE_PARKING_KIND
	default:
		return sso.ParkingKind_UNDEFINED_PARKING_KIND
	}
}
