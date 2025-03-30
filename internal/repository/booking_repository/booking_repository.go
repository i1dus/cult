package booking_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
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

	var rentalID uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT id FROM rentals 
         WHERE parking_lot_id = $1 
         AND start_at <= $2 
         AND end_at >= $3`,
		booking.ParkingLot,
		booking.From,
		booking.To,
	).Scan(&rentalID)

	query := `
        INSERT INTO bookings (rental_id, user_id, vehicle, start_at, end_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (rental_id, start_at, end_at) 
        DO NOTHING
    `
	res, err := r.db.Exec(
		ctx,
		query,
		rentalID,
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

func (r *BookingRepository) EditBooking(ctx context.Context, bookingID uuid.UUID, to time.Time) error {
	const op = "BookingRepository.EditBooking"

	res, err := r.db.Exec(ctx,
		`UPDATE bookings SET end_at = $1 WHERE id = $2`,
		to,
		bookingID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrBookingNotFound)
	}

	return nil
}

// GetBooking implements
func (r *BookingRepository) GetBooking(ctx context.Context, parkingLot int64) (*domain.Booking, error) {
	const op = "BookingRepository.GetBooking"

	var rentalID uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT id FROM rentals 
         WHERE parking_lot_id = $1 
         AND NOW() BETWEEN start_at AND end_at`,
		parkingLot,
	).Scan(&rentalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `
        SELECT user_id, vehicle, is_short_term, is_present, start_at, end_at
        FROM bookings
        WHERE rental_id = $1
          AND is_present
    `

	var booking domain.Booking
	err = r.db.QueryRow(ctx, query, rentalID).Scan(
		&booking.UserID,
		&booking.Vehicle,
		&booking.IsShortTerm,
		&booking.IsPresent,
		&booking.From,
		&booking.To,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	booking.ParkingLot = parkingLot
	booking.RentalID = rentalID

	return &booking, nil
}

// GetBookingsByFilter implements
func (r *BookingRepository) GetBookingsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Booking, error) {
	const op = "BookingRepository.GetBookingsByFilter"

	query := `
        SELECT r.parking_lot_id, b.user_id, b.vehicle, b.start_at, b.end_at
        FROM bookings b
        JOIN rentals r ON b.rental_id = r.id
        WHERE ($1::uuid IS NULL OR b.user_id = $1)
          AND (
              ($2::timestamp IS NULL AND $3::timestamp IS NULL) OR
              (b.start_at <= $3 AND b.end_at >= $2)
          )
          AND ($4::int[] IS NULL OR r.parking_lot_id = ANY($4::int[]))
        ORDER BY b.start_at
    `

	// Convert filter parameters to SQL-compatible values
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
		return nil, fmt.Errorf("%s: querying bookings: %w", op, err)
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
			return nil, fmt.Errorf("%s: scanning booking: %w", op, err)
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return bookings, nil
}

// GetParkingLotsByFilter implements
func (r *BookingRepository) GetParkingLotsByFilter(ctx context.Context, filter domain.Filter) ([]domain.ParkingLot, error) {
	const op = "BookingRepository.GetParkingLotsByFilter"

	query := `
        SELECT p.id, p.owner_vehicle
        FROM parking_lots p
        JOIN rentals r ON p.id = r.parking_lot_id
        WHERE NOT EXISTS (
            SELECT 1 FROM rentals r
            WHERE r.parking_lot_id = p.id
            AND (
                ($1 BETWEEN r.start_at AND r.end_at) OR  -- New booking starts during existing rental
                ($2 BETWEEN r.start_at AND r.end_at) OR  -- New booking ends during existing rental
                (r.start_at BETWEEN $1 AND $2) OR        -- Existing rental starts during new booking
                (r.end_at BETWEEN $1 AND $2)             -- Existing rental ends during new booking
            )
        )`

	args := []interface{}{filter.From, filter.To}

	// Add optional filters
	if filter.UserID != uuid.Nil {
		query += " AND p.owner_id = $3"
		args = append(args, filter.UserID)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query available parking lots: %w", err)
	}
	defer rows.Close()

	var parkingLots []domain.ParkingLot
	for rows.Next() {
		var lot domain.ParkingLot
		err := rows.Scan(
			&lot.ID,
			&lot.OwnerVehicle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parking lot: %w", err)
		}
		parkingLots = append(parkingLots, lot)
	}

	return parkingLots, nil
}
