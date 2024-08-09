package usecase

import (
	"avito_bootcamp/internal/entity"
	"context"
	"fmt"
)

type FlatUseCase struct {
	repo FlatRepo
}

var _ Flat = (*FlatUseCase)(nil)

func NewFlatUseCase(r FlatRepo) *FlatUseCase {
	return &FlatUseCase{
		repo: r,
	}
}

func (uc *FlatUseCase) Create(ctx context.Context, f entity.Flat) (entity.Flat, error) {
	flat, err := uc.repo.Store(context.Background(), f)
	if err != nil {
		return entity.Flat{}, fmt.Errorf("FlatUseCase - Create - s.repo.Store: %w", err)
	}

	return flat, nil
}

func (uc *FlatUseCase) Update(ctx context.Context, f entity.Flat) (entity.Flat, error) {
	flat, err := uc.repo.Update(context.Background(), f)
	if err != nil {
		return entity.Flat{}, fmt.Errorf("FlatUseCase - Update - s.repo.Update: %w", err)
	}

	return flat, nil
}
