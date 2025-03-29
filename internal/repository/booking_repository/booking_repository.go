package booking_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	desc "cult/pkg"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/protobuf/types/known/timestamppb"

    "github.com/lib/pq"
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


// GetBookingsByFilter implements
func (r *BookingRepository) GetBookingsByFilter(ctx context.Context, db *sql.DB, filter domain.Filter) ([]domain.Booking, error) {
	query := `
        SELECT b.parking_lot_id, b.user_id, b.vehicle_id, b.start_at, b.end_at
        FROM bookings b
        JOIN parking_lots pl ON b.parking_lot_id = pl.id::TEXT
        WHERE pl.owner_id = $1 OR $2 IS NULL
          AND (
              ($3::timestamp IS NULL AND $4::timestamp IS NULL) OR
              (b.start_at <= $4 AND b.end_at >= $3)
          )
          AND ($5::int[] IS NULL OR b.parking_lot_id::int = ANY($5::int[]))
        ORDER BY b.start_at
    `

	var from, to interface{}
	if filter.From.IsZero() {
		from = nil
	} else {
		from = filter.From
	}

	if filter.To.IsZero() {
		to = nil
	} else {
		to = filter.To
	}

	var parkingLots interface{}
	if len(filter.ParkingLots) == 0 {
		parkingLots = nil
	} else {
		parkingLots = pq.Array(filter.ParkingLots)
	}

	rows, err := db.QueryContext(ctx, query,
		filter.UserID,
		from,
		to,
		parkingLots,
	)

	if err != nil {
		return nil, fmt.Errorf("querying bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(
			&b.ParkingLot,
			&b.UserID,
			&b.Vehicle,
			&b.From,
			&b.To,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning booking: %w", err)
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
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
