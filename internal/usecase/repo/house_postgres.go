package repo

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/postgres"
	"context"
	"fmt"
	"time"
)

var _ usecase.HouseRepo = (*HouseRepo)(nil)

type HouseRepo struct {
	pg *postgres.Postgres
}

func NewHouseRepo(pg *postgres.Postgres) *HouseRepo {
	return &HouseRepo{
		pg: pg,
	}
}

func (r *HouseRepo) Store(ctx context.Context, h entity.House) (entity.House, error) {
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	query := `INSERT INTO houses (id, address, year, developer, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pg.Pool.Exec(ctx, query, h.ID, h.Address, h.Year, h.Developer, h.CreatedAt, h.UpdatedAt)

	if err != nil {
		return entity.House{}, fmt.Errorf("HouseRepo - Store - r.Pool.Exec: %w", err)
	}

	return h, nil

}

func (r *HouseRepo) GetFlats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error) {
	var args []any
	query := "SELECT id, number, house_id, price, rooms, status FROM flats WHERE house_id = $1"
	args = append(args, houseID)

	if userType == "client" {
		query += " AND status = $2"
		args = append(args, "approved")
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("HouseRepo - GetFlatsByHouseID - r.pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var flats []entity.Flat
	for rows.Next() {
		var flat entity.Flat
		if err := rows.Scan(&flat.ID, &flat.Number, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return nil, fmt.Errorf("HouseRepo - GetFlatsByHouseID - rows.Scan: %w", err)
		}
		flats = append(flats, flat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("HouseRepo - GetFlatsByHouseID - rows.Err: %w", err)
	}

	return flats, nil
}
