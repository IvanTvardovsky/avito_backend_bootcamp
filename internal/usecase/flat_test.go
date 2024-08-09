package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
)

type MockFlatRepo struct {
	mock.Mock
}

func (m *MockFlatRepo) Store(ctx context.Context, f entity.Flat) (entity.Flat, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(entity.Flat), args.Error(1)
}

func (m *MockFlatRepo) Update(ctx context.Context, f entity.Flat) (entity.Flat, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(entity.Flat), args.Error(1)
}

func (m *MockFlatRepo) GetFlats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error) {
	args := m.Called(ctx, houseID, userType)
	return args.Get(0).([]entity.Flat), args.Error(1)
}

func TestFlatUseCase_Create_Success(t *testing.T) {
	mockRepo := new(MockFlatRepo)
	uc := usecase.NewFlatUseCase(mockRepo)

	flat := entity.Flat{
		Number:  101,
		HouseID: 1,
		Price:   100000,
		Rooms:   3,
		Status:  "pending",
	}

	expectedFlat := flat
	expectedFlat.ID = 1

	mockRepo.On("Store", mock.Anything, flat).Return(expectedFlat, nil)

	createdFlat, err := uc.Create(context.Background(), flat)

	assert.NoError(t, err)
	assert.Equal(t, expectedFlat, createdFlat)
	mockRepo.AssertExpectations(t)
}

func TestFlatUseCase_Create_StoreError(t *testing.T) {
	mockRepo := new(MockFlatRepo)
	uc := usecase.NewFlatUseCase(mockRepo)

	flat := entity.Flat{
		Number:  101,
		HouseID: 1,
		Price:   100000,
		Rooms:   3,
		Status:  "pending",
	}

	expectedError := errors.New("some store error")

	mockRepo.On("Store", mock.Anything, flat).Return(entity.Flat{}, expectedError)

	_, err := uc.Create(context.Background(), flat)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, expectedError))
	mockRepo.AssertExpectations(t)
}
