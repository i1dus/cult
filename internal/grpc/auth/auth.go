package authgrpc

import (
	"context"
	"cult/internal/domain"
	desc "cult/pkg"

	"github.com/google/uuid"

	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context, phoneNumber string, password string) (uuid.UUID, string, error)
	RegisterNewUser(ctx context.Context, phoneNumber string, password string) (userID uuid.UUID, err error)
	GetUserByPhone(ctx context.Context, phoneNumber string) (*domain.User, error)
}

type serverAPI struct {
	desc.UnimplementedParkingAPIServer

	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	desc.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth})
}
