package rental_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"math"
	"time"

	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type RentalRepository struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewRentalRepository(db *pgx.Conn, log *slog.Logger) *RentalRepository {
	return &RentalRepository{
		db:  db,
		log: log,
	}
}

// AddRental implements add interface
func (r *RentalRepository) AddRental(ctx context.Context, rental domain.Rental) error {
	const op = "RentalRepository.AddRental"

	if rental.To.Before(rental.From) {
		return fmt.Errorf("%s: %w", op, repository.ErrInvalidTimeRange)
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(
            SELECT 1 FROM rentals 
            WHERE parking_lot_id = $1 
            AND (
                ($2 BETWEEN start_at AND end_at) OR 
                ($3 BETWEEN start_at AND end_at) OR 
                (start_at BETWEEN $2 AND $3) OR 
                (end_at BETWEEN $2 AND $3)
            )
        )`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check existing rentals: %w", err)
	}

	if exists {
		return repository.ErrRentalConflict
	}

	_, err = tx.Exec(ctx, `INSERT INTO rentals (
            parking_lot_id,
            start_at,
            end_at,
            cost_per_hour,
            cost_per_day
        ) VALUES ($1, $2, $3, $4, $5)`,
		rental.ParkingLot,
		rental.From,
		rental.To,
		rental.CostPerHour,
		rental.CostPerDay,
	)
	if err != nil {
		return fmt.Errorf("failed to add rental: %w", err)
	}

	return tx.Commit(ctx)
}

// GetRentalsByFilter implements
func (r *RentalRepository) GetRentalsByFilter(ctx context.Context, filter domain.Filter) ([]domain.Rental, error) {
	const op = "RentalRepository.GetRentalsByFilter"

	parkingLotRows, err := r.db.Query(ctx, `
        SELECT id FROM parking_lots WHERE parking_lots.owner_id = $1`,
		filter.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's parking lots: %w", err)
	}
	defer parkingLotRows.Close()

	var parkingLotIDs []int64
	for parkingLotRows.Next() {
		var id int64
		if err := parkingLotRows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan parking lot ID: %w", err)
		}
		parkingLotIDs = append(parkingLotIDs, id)
	}

	if len(parkingLotIDs) == 0 {
		return []domain.Rental{}, nil
	}

	query := `
        SELECT 
            r.id,
            r.parking_lot_id,
            r.start_at,
            r.end_at,
            r.cost_per_hour,
            r.cost_per_day
        FROM rentals r
        WHERE r.parking_lot_id = ANY($1)`

	args := []interface{}{pq.Array(parkingLotIDs)}
	argPos := 2

	if !filter.From.IsZero() {
		query += fmt.Sprintf(" AND r.start_at >= $%d", argPos)
		args = append(args, filter.From)
		argPos++
	}

	if !filter.To.IsZero() {
		query += fmt.Sprintf(" AND r.end_at <= $%d", argPos)
		args = append(args, filter.To)
		argPos++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query rentals: %w", err)
	}
	defer rows.Close()

	var rentals []domain.Rental
	for rows.Next() {
		var rental domain.Rental
		if err := rows.Scan(
			&rental.ID,
			&rental.ParkingLot,
			&rental.From,
			&rental.To,
			&rental.CostPerHour,
			&rental.CostPerDay,
		); err != nil {
			return nil, fmt.Errorf("failed to scan rental: %w", err)
		}
		rentals = append(rentals, domain.Rental{})
	}

	return rentals, nil
}

// GetBookingPriceByID implements
func (r *RentalRepository) GetBookingPriceByID(ctx context.Context, id uuid.UUID) (int64, error) {
	const op = "RentalRepository.GetBookingPriceByID"
	price := int64(0)

	rows, err := r.db.Query(ctx, `
        SELECT rental_id, is_short_term FROM bookings WHERE id = $1`,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get booking price: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("no booking found with id %s", id)
	}

	var rentalID uuid.UUID
	var isShortTerm bool

	err = rows.Scan(&rentalID, &isShortTerm)

	query := `
        SELECT 
            r.id,
            r.parking_lot_id,
            r.start_at,
            r.end_at,
            r.cost_per_hour,
            r.cost_per_day
        FROM rentals r
        WHERE r.id = $1`

	var rental domain.Rental

	resRow := r.db.QueryRow(ctx, query, rentalID).Scan(&rental)
	if err != nil {
		return nil, fmt.Errorf("failed to query rentals: %w", err)
	}
	
	domain.Rental{
		ID:          rentalID,
		ParkingLot:  0,
		From:        time.Time{},
		To:          time.Time{},
		CostPerHour: 0,
		CostPerDay:  0,
	}

	duration := rental.To.Sub(rental.From)



	if isShortTerm {
		hours := int(math.Ceil(duration.Minutes()) * 10)

		price = rental.CostPerHour * hours
	} else {
		days := int(math.Ceil(duration.Hours()) / 24)
		price = rental.CostPerDay *
	}
	
	return price, nil
}
