package grpc

import (
	"context"
	sso "cult/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Login(ctx context.Context, in *sso.LoginRequest) (*sso.LoginResponse, error) {
	if in.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, token, err := s.auth.Login(ctx, in.GetPhoneNumber(), in.GetPassword())
	if err != nil {
		//if errors.Is(err, ErrInvalidCredentials) {
		//	return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		//}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.LoginResponse{Token: token, UserID: userID.String()}, nil
}
