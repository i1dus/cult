package grpc

import (
	"context"
	"cult/internal/domain"
	sso "cult/pkg"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) GetUserByPhoneNumber(ctx context.Context, in *sso.GetUserByPhoneRequest) (*sso.GetUserByPhoneResponse, error) {
	if in.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}

	user, err := s.auth.GetUserByPhone(ctx, in.GetPhoneNumber())
	if err != nil {
		//if errors.Is(err, ErrInvalidCredentials) {
		//	return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		//}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.GetUserByPhoneResponse{
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
