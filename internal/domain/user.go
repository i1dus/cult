package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID
	Phone      string
	Name       string
	Surname    string
	Patronymic string
	Address    string
	UserType   UserType
	PassHash   []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserType string

const (
	UndefinedUserType       UserType = "UNDEFINED_USER_TYPE"
	RegularUserType         UserType = "REGULAR_USER_TYPE"
	ManagingCompanyUserType UserType = "MANAGING_COMPANY_USER_TYPE"
	AdministratorUserType   UserType = "ADMINISTRATOR_USER_TYPE"
)

func (u UserType) String() string {
	return string(u)
}
