package grpcapp

import (
	"context"
	"crypto/tls"
	authgrpc "cult/internal/grpc"
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
}

func New(
	log *slog.Logger,
	authService authgrpc.AuthService,
	parkingLotService authgrpc.ParkingLotService,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
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
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 2)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

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

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcGatewayMux)

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
	swaggerMux.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	httpMux.Handle("/swagger/", swaggerMux)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      httpMux,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

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

	a.gRPCServer.GracefulStop()
	return shutdownErr
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
