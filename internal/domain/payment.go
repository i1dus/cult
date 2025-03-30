package domain

import (
	"time"
)

type Payment struct {
	ID          string
	UserID      string
	BookingID   string
	RentalID    string
	PaymentType PaymentType
	Amount      int64
	Currency    string
	Status      PaymentStatus
	Method      PaymentMethod
	CreatedAt   time.Time
	PaidAt      time.Time
	RefundedAt  time.Time
}

type PaymentType int32

const (
	PaymentType_UNKNOWN PaymentType = iota
	PaymentType_BOOKING
	PaymentType_RENTAL
)

type PaymentStatus int32

const (
	PaymentStatus_UNDEFINED_PAYMENT_STATUS PaymentStatus = 0
	PaymentStatus_PENDING                  PaymentStatus = 1
	PaymentStatus_PROCESSING               PaymentStatus = 2
	PaymentStatus_COMPLETED                PaymentStatus = 3
	PaymentStatus_FAILED                   PaymentStatus = 4
	PaymentStatus_REFUNDED                 PaymentStatus = 5
	PaymentStatus_PARTIALLY_REFUNDED       PaymentStatus = 6
	PaymentStatus_CANCELLED                PaymentStatus = 7
)

type PaymentMethod int32

const (
	PaymentMethod_UNDEFINED_PAYMENT_METHOD PaymentMethod = 0
	PaymentMethod_CREDIT_CARD              PaymentMethod = 1
	PaymentMethod_BANK_TRANSFER            PaymentMethod = 2
	PaymentMethod_MOBILE_PAYMENT           PaymentMethod = 3
	PaymentMethod_ELECTRONIC_WALLET        PaymentMethod = 4
)
