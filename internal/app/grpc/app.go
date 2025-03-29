package grpcapp

import (
	"context"
	"crypto/tls"
	authgrpc "cult/internal/grpc/auth"
	api "cult/pkg"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc/credentials/insecure"

	httpSwagger "github.com/swaggo/http-swagger"

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
	authService authgrpc.AuthService,
	parkingLotService authgrpc.ParkingLotService,
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

	authgrpc.Register(gRPCServer, authService, parkingLotService)

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

	ctx := context.Background()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Error channel to collect server errors
	errChan := make(chan error, 3)

	// Create TCP listener for gRPC
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Create HTTP router for gRPC gateway
	grpcGatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = api.RegisterParkingAPIHandlerFromEndpoint(
		ctx,
		grpcGatewayMux,
		fmt.Sprintf("localhost:%d", a.port),
		opts,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Main HTTP server
	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcGatewayMux)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      httpMux,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	// Swagger server
	swaggerMux := chi.NewMux()
	swaggerMux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("pkg/parking.swagger.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	swaggerMux.HandleFunc("/swagger/*", httpSwagger.WrapHandler)
	swaggerServer := &http.Server{
		Addr:    ":8081", // Changed to standard port
		Handler: swaggerMux,
	}

	// Start servers in goroutines
	go func() {
		a.log.Info("gRPC server starting", slog.String("addr", grpcListener.Addr().String()))
		if err := a.gRPCServer.Serve(grpcListener); err != nil && err != grpc.ErrServerStopped {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		a.log.Info("HTTP server starting", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	go func() {
		a.log.Info("Swagger UI available", slog.String("addr", swaggerServer.Addr))
		if err := swaggerServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("Swagger server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("shutdown signal received")
	case err := <-errChan:
		cancel()
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("shutting down servers")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	var shutdownErr error

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		shutdownErr = fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	if err := swaggerServer.Shutdown(shutdownCtx); err != nil {
		shutdownErr = fmt.Errorf("Swagger server shutdown error: %w", err)
	}

	a.gRPCServer.GracefulStop() // Graceful stop for gRPC

	return shutdownErr
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
