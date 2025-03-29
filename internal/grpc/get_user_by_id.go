package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"
	"fmt"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) GetUserByID(ctx context.Context, in *sso.GetUserByIDRequest) (*sso.GetUserByIDResponse, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	parsedUserID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot convert user ID to UUID")
	}
	user, err := s.auth.GetUserByID(ctx, parsedUserID)
	if err != nil {
		//if errors.Is(err, ErrInvalidCredentials) {
		//	return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		//}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user: %s", err.Error()))
	}

	return &sso.GetUserByIDResponse{
		User: &sso.User{
			Id:          user.ID.String(),
			Name:        user.Name,
			Surname:     user.Surname,
			Patronymic:  user.Patronymic,
			PhoneNumber: user.Phone,
			Address:     user.Address,
			UserType:    userType(user.UserType),
		},
	}, nil
}

func userType(uType domain.UserType) sso.UserType {
	switch uType {
	case domain.RegularUserType:
		return sso.UserType_REGULAR_USER_TYPE
	case domain.ManagingCompanyUserType:
		return sso.UserType_MANAGING_COMPANY_USER_TYPE
	case domain.AdministratorUserType:
		return sso.UserType_ADMINISTRATOR_USER_TYPE
	}
	return sso.UserType_UNDEFINED_USER_TYPE
}
