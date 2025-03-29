package authgrpc

import (
	"context"
	"cult/internal/repository"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

import "cult/pkg"

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	in *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &sso.IsAdminResponse{IsAdmin: isAdmin}, nil
}
