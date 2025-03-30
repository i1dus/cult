package grpc

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	sso "cult/pkg"
	"errors"
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

func (s *serverAPI) UpdateUser(ctx context.Context, req *sso.UpdateUserRequest) (*sso.UpdateUserResponse, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	update := domain.UserUpdate{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Phone:      req.PhoneNumber,
		Address:    req.Address,
	}

	if err := s.auth.UpdateUser(ctx, userID, update); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "phone number already in use")
		}
		if errors.Is(err, repository.ErrNoFieldsToUpdate) {
			return nil, status.Error(codes.InvalidArgument, "no fields to update")
		}
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &sso.UpdateUserResponse{}, nil
}

func (s *serverAPI) GetUserByPhoneNumber(ctx context.Context, in *sso.GetUserByPhoneNumberRequest) (*sso.GetUserByPhoneNumberResponse, error) {
	if in.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	user, err := s.auth.UserByPhoneNumber(ctx, in.PhoneNumber)
	if err != nil {
		//if errors.Is(err, ErrInvalidCredentials) {
		//	return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		//}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user: %s", err.Error()))
	}

	return &sso.GetUserByPhoneNumberResponse{
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
