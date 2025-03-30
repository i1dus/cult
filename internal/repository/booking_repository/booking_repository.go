package booking_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
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
func (r *BookingRepository) AddBooking(ctx context.Context, booking domain.Booking) error {
	const op = "BookingRepository.AddBooking"

	if booking.To.Before(booking.From) {
		return fmt.Errorf("%s: %w", op, repository.ErrInvalidTimeRange)
	}

	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM parking_lots WHERE id = $1)",
		booking.ParkingLot,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, "parking_lots does not exist")
	}

	query := `
        INSERT INTO bookings (parking_lot_id, user_id, vehicle_id, start_at, end_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (parking_lot_id, start_at, end_at) 
        DO NOTHING
    `
	res, err := r.db.Exec(
		ctx,
		query,
		booking.ParkingLot,
		booking.UserID,
		booking.Vehicle,
		booking.From,
		booking.To,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrBookingConflict)
	}

	return nil
}

// GetBooking implements
func (r *BookingRepository) GetBooking(ctx context.Context, parkingLot int64) (domain.Booking, error) {
	const op = "BookingRepository.GetBooking"

	query := `
		SELECT user_id, vehicle_id, start_at, end_at
		FROM bookings
		WHERE parking_lot_id = $1
		  AND NOW() BETWEEN start_at AND end_at
	`

	var booking domain.Booking
	err := r.db.QueryRow(ctx, query, parkingLot).Scan(
		&booking.UserID,
		&booking.Vehicle,
		&booking.From,
		&booking.To,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf("%s: %w", op, repository.ErrBookingNotFound)
		}
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}

// GetBookingsByFilter implements
func (r *BookingRepository) GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error) {
	query := `
        SELECT b.parking_lot_id, b.user_id, b.vehicle_id, b.start_at, b.end_at
        FROM bookings b
        JOIN parking_lots pl ON b.parking_lot_id = pl.id
        WHERE pl.owner_id = $1 OR $1 IS NULL
          AND (
              ($2::timestamp IS NULL AND $3::timestamp IS NULL) OR
              (b.start_at <= $3 AND b.end_at >= $2)
          )
          AND ($4::int[] IS NULL OR b.parking_lot_id::int = ANY($4::int[]))
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

	rows, err := r.db.Query(ctx, query,
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
