package usecase

import (
	"avito_bootcamp/internal/entity"
	"context"
	"fmt"
)

type HouseUseCase struct {
	repo HouseRepo
}

var _ House = (*HouseUseCase)(nil)

func NewHouseUseCase(r HouseRepo) *HouseUseCase {
	return &HouseUseCase{
		repo: r,
	}
}

func (uc *HouseUseCase) Create(ctx context.Context, h entity.House) (entity.House, error) {
	house, err := uc.repo.Store(context.Background(), h)
	if err != nil {
		return entity.House{}, fmt.Errorf("HouseUseCase - Create - s.repo.Store: %w", err)
	}

	return house, nil
}

func (uc *HouseUseCase) Flats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error) {
	flats, err := uc.repo.GetFlats(ctx, houseID, userType)
	if err != nil {
		return nil, fmt.Errorf("HouseUseCase - GetFlatsByHouseID - uc.repo.GetFlatsByHouseID: %w", err)
	}

	return flats, nil
}
