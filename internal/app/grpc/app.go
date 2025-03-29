package grpcapp

import (
	"context"
	"crypto/tls"
	authgrpc "cult/internal/grpc/auth"
	api "cult/pkg"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	host       string
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		host:       "localhost:44044",
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	op := "app.Run"

	// Create a TCP listener for gRPC
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Create HTTP router for gRPC gateway
	grpcGatewayMux := runtime.NewServeMux()

	// Register gRPC gateway endpoints
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = api.RegisterParkingAPIHandlerFromEndpoint(
		context.Background(),
		grpcGatewayMux,
		fmt.Sprintf("localhost:%d", a.port),
		opts,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcGatewayMux)

	// HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: httpMux,
		// Enforce HTTP/1.1
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	// Use WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Start gRPC server in goroutine
	go func() {
		defer wg.Done()
		a.log.Info("gRPC server started", slog.String("addr", grpcListener.Addr().String()))
		if err := a.gRPCServer.Serve(grpcListener); err != nil {
			a.log.Error("gRPC server failed", slog.String("error", err.Error()))
		}
	}()

	// Start HTTP server in goroutine
	go func() {
		defer wg.Done()
		a.log.Info("HTTP server started", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	// Wait for servers to exit
	wg.Wait()
	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
