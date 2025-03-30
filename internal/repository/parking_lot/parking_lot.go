package parking_lot

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strings"
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
		SELECT id, parking_kind, owner_id, owner_vehicle
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

		var parkingLotKind string
		err := rows.Scan(
			&lot.ID,
			&parkingLotKind,
			&lot.OwnerID,
			&lot.OwnerVehicle,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lot.ParkingKind = domain.ParkingKind(parkingLotKind)
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
        SELECT id, parking_kind, owner_id, owner_vehicle
        FROM parking_lots
        WHERE id = $1
    `

	var lot domain.ParkingLot
	var parkingLotKind string

	err := r.db.QueryRow(ctx, query, number).Scan(
		&lot.ID,
		&parkingLotKind,
		&lot.OwnerID,
		&lot.OwnerVehicle,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return domain.ParkingLot{}, fmt.Errorf("%s: %w", op, err)
	}

	lot.ParkingKind = domain.ParkingKind(parkingLotKind)

	return lot, nil
}

func (r *ParkingLotRepo) GetParkingLotsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.ParkingLot, error) {
	const op = "parkingLotRepo.GetParkingLotsByOwnerID"

	query := `
        SELECT id, parking_kind, owner_id, owner_vehicle
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
		var parkingLotKind string

		err := rows.Scan(
			&lot.ID,
			&parkingLotKind,
			&lot.OwnerID,
			&lot.OwnerVehicle,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		lot.ParkingKind = domain.ParkingKind(parkingLotKind)

		lots = append(lots, lot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lots, nil
}

func (r *ParkingLotRepo) UpdateParkingLot(ctx context.Context, parkingLotID string, update domain.ParkingLotUpdate) error {
	const op = "ParkingLotRepo.UpdateParkingLot"

	// Build dynamic update query
	query := "UPDATE parking_lots SET "
	params := []interface{}{}
	setClauses := []string{}
	paramCount := 1

	if update.ParkingKind != nil {
		setClauses = append(setClauses, fmt.Sprintf("parking_kind = $%d", paramCount))
		params = append(params, update.ParkingKind.String())
		paramCount++
	}

	if update.OwnerID != nil {
		setClauses = append(setClauses, fmt.Sprintf("owner_id = $%d", paramCount))
		params = append(params, *update.OwnerID)
		paramCount++
	}

	if update.OwnerVehicle != nil {
		setClauses = append(setClauses, fmt.Sprintf("owner_vehicle = $%d", paramCount))
		params = append(params, *update.OwnerVehicle)
		paramCount++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNoFieldsToUpdate)
	}

	// Add WHERE clause
	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, parkingLotID)

	// Execute query
	tag, err := r.db.Exec(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrParkingLotNotFound)
	}

	return nil
}
