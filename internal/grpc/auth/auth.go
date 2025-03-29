package authgrpc

import (
	"context"
	"cult/internal/gen/parking_lot"
	"cult/internal/gen/sso"

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
	sso.UnimplementedAuthServer
	parking_lot.UnimplementedParkingAPIServer

	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
	parking_lot.RegisterParkingAPIServer(gRPCServer, &serverAPI{auth: auth})
}
