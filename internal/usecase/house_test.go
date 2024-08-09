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

type MockHouseRepo struct {
	mock.Mock
}

func (m *MockHouseRepo) GetFlats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error) {
	args := m.Called(ctx, houseID, userType)
	return args.Get(0).([]entity.Flat), args.Error(1)
}

func (m *MockHouseRepo) Store(ctx context.Context, h entity.House) (entity.House, error) {
	args := m.Called(ctx, h)
	return args.Get(0).(entity.House), args.Error(1)
}

func TestHouseUseCase_Flats_Success(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	uc := usecase.NewHouseUseCase(mockRepo)

	houseID := 1
	userType := "client"
	expectedFlats := []entity.Flat{
		{ID: 1, Number: 101, HouseID: houseID, Price: 100000, Rooms: 3, Status: "approved"},
		{ID: 2, Number: 102, HouseID: houseID, Price: 120000, Rooms: 4, Status: "approved"},
	}

	mockRepo.On("GetFlats", mock.Anything, houseID, userType).Return(expectedFlats, nil)

	flats, err := uc.Flats(context.Background(), houseID, userType)

	assert.NoError(t, err)
	assert.Equal(t, expectedFlats, flats)
	mockRepo.AssertExpectations(t)
}

func TestHouseUseCase_Flats_Error(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	uc := usecase.NewHouseUseCase(mockRepo)

	houseID := 1
	userType := "client"
	expectedError := errors.New("some error")

	mockRepo.On("GetFlats", mock.Anything, houseID, userType).Return([]entity.Flat(nil), expectedError)

	flats, err := uc.Flats(context.Background(), houseID, userType)

	assert.Error(t, err)
	assert.Nil(t, flats)
	assert.True(t, errors.Is(err, expectedError))
	mockRepo.AssertExpectations(t)
}
