package parking_lot

import (
	"context"
	"cult/internal/domain"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type ParkingLotRepository interface {
	GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error)
}

type parkingLotRepo struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewParkingLotRepository(db *pgx.Conn, log *slog.Logger) ParkingLotRepository {
	return &parkingLotRepo{
		db:  db,
		log: log,
	}
}

func (r *parkingLotRepo) GetAllParkingLots(ctx context.Context) ([]domain.ParkingLot, error) {
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
