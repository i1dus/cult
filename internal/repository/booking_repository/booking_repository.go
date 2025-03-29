package booking_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	desc "cult/pkg"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jackc/pgx/v5"
)

type BookingRepository struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewBookingRepository(db *pgx.Conn, log *slog.Logger) *BookingRepository {
	return &BookingRepository{
		db:  db,
		log: log,
	}
}

// AddBooking implements UserSaver interface
func (r *BookingRepository) AddBooking(ctx context.Context, parkingLot, carPlate string, userId int64, from, to timestamppb.Timestamp) error {
	const op = "BookingRepository.AddBooking"

	query := `
		INSERT INTO bookings (parking_lot_id, user_id, vehicle_id, start_at, end_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	err := r.db.QueryRow(ctx, query, parkingLot, userId, carPlate, userId, from, to).Scan()
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%s: %w", op, repository.ErrBookingExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetBooking implements
func (r *BookingRepository) GetBooking(ctx context.Context, parkingLot string) error {
	const op = "BookingRepository.GetBooking"

	desc.GetParkingBookingRequest{ParkingLot: desc.ParkingLot{
		Number:  0,
		Type:    0,
		Status:  0,
		Vehicle: nil,
		Owner:   nil,
	}}

	query := `
		SELECT user_id, vehicle_id, start_at, end_at
		FROM bookings
		WHERE parking_lot_id = $1
		  AND NOW() BETWEEN start_at AND end_at
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, parkingLot).Scan(
		&user.ID,
		&user.,
		&user.PassHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// IsAdmin implements UserProvider interface
func (r *BookingRepository) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "UserRepository.IsAdmin"

	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE id = $1 AND user_type = 'admin'
		)
	`

	var isAdmin bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

// helper function to check for unique violation
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
