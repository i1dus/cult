package authgrpc

import (
	"context"
	"cult/internal/repository"
	sso "cult/pkg"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Register(ctx context.Context, in *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if in.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.GetPhoneNumber(), in.GetPassword())
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &sso.RegisterResponse{UserId: userID.String()}, nil
}
