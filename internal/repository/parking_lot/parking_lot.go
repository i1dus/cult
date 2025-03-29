package parking_lot

import (
	"context"
	"cult/internal/domain"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

var ErrNotFound = errors.New("parking lot not found")

type ParkingLotRepo struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewParkingLotRepository(db *pgx.Conn, log *slog.Logger) *ParkingLotRepo {
	return &ParkingLotRepo{
		db:  db,
		log: log,
	}
}

func (r *ParkingLotRepo) GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error) {
	const op = "parkingLotRepo.GetAllParkingLots"

	query := `
		SELECT id, parking_type, vehicle_id, owner_id
		FROM parking_lots
		ORDER BY id::integer
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var lots []domain.ParkingLot
	for rows.Next() {
		var lot domain.ParkingLot

		var parkingLotType string
		err := rows.Scan(
			&lot.ID,
			&parkingLotType,
			&lot.OwnerID,
			&lot.VehicleID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lot.ParkingType = domain.ParkingType(parkingLotType)

		lots = append(lots, lot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lots, nil
}

func (r *ParkingLotRepo) GetParkingLotByNumber(ctx context.Context, number string) (domain.ParkingLot, error) {
	const op = "parkingLotRepo.GetParkingLotByNumber"

	query := `
        SELECT id, parking_type, owner_id
        FROM parking_lots
        WHERE id = $1
    `

	var lot domain.ParkingLot
	var parkingLotType string

	err := r.db.QueryRow(ctx, query, number).Scan(
		&lot.ID,
		&parkingLotType,
		&lot.OwnerID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, err)
	}

	lot.ParkingType = domain.ParkingType(parkingLotType)

	return lot, nil
}

func (r *ParkingLotRepo) GetParkingLotsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error) {
	const op = "parkingLotRepo.GetParkingLotsByOwnerID"

	query := `
        SELECT id, parking_type, owner_id
        FROM parking_lots
        WHERE owner_id = $1
        ORDER BY id::integer
    `

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var lots []domain.ParkingLot
	for rows.Next() {
		var lot domain.ParkingLot
		var parkingLotType string

		err := rows.Scan(
			&lot.ID,
			&parkingLotType,
			&lot.OwnerID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lot.ParkingType = domain.ParkingType(parkingLotType)

		lots = append(lots, lot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lots, nil
}
