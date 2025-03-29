package authgrpc

import (
	"context"
	desc "cult/pkg"

	"google.golang.org/grpc"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	desc.UnimplementedParkingAPIServer

	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	desc.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth})
}
