package repo

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/postgres"
	"context"
	"fmt"
	"strings"
)

var _ usecase.FlatRepo = (*FlatRepo)(nil)

type FlatRepo struct {
	pg *postgres.Postgres
}

func NewFlatRepo(pg *postgres.Postgres) *FlatRepo {
	return &FlatRepo{
		pg: pg,
	}
}

func (r *FlatRepo) Store(ctx context.Context, flat entity.Flat) (entity.Flat, error) {
	flat.Status = "created"

	query := `INSERT INTO flats (number, house_id, price, rooms, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.pg.Pool.QueryRow(ctx, query, flat.Number, flat.HouseID, flat.Price, flat.Rooms, flat.Status).Scan(&flat.ID)

	if err != nil {
		return entity.Flat{}, fmt.Errorf("FlatRepo - Store - r.Pool.Exec: %w", err)
	}

	return flat, nil
}

func (r *FlatRepo) Update(ctx context.Context, f entity.Flat) (entity.Flat, error) {
	var updates []string
	var args []interface{}
	idx := 1

	if f.Number != -1 {
		updates = append(updates, fmt.Sprintf("number = $%d", idx))
		args = append(args, f.Number)
		idx++
	}

	if f.HouseID != -1 {
		updates = append(updates, fmt.Sprintf("house_id = $%d", idx))
		args = append(args, f.HouseID)
		idx++
	}

	if f.Price != -1 {
		updates = append(updates, fmt.Sprintf("price = $%d", idx))
		args = append(args, f.Price)
		idx++
	}

	if f.Rooms != -1 {
		updates = append(updates, fmt.Sprintf("rooms = $%d", idx))
		args = append(args, f.Rooms)
		idx++
	}

	if f.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}

	if len(updates) == 0 {
		return entity.Flat{}, fmt.Errorf("FlatRepo - Update - nothing to update")
	}

	setClause := strings.Join(updates, ", ")

	query := fmt.Sprintf("UPDATE flats SET %s WHERE id = $%d RETURNING id, number, house_id, price, rooms, status", setClause, idx)
	args = append(args, f.ID)

	var updatedFlat entity.Flat

	err := r.pg.Pool.QueryRow(ctx, query, args...).Scan(&updatedFlat.ID, &updatedFlat.Number, &updatedFlat.HouseID, &updatedFlat.Price, &updatedFlat.Rooms, &updatedFlat.Status)
	if err != nil {
		return f, err
	}

	return updatedFlat, nil
}
