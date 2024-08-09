package usecase

import (
	"avito_bootcamp/internal/entity"
	"context"
)

type (
	Flat interface {
		Create(ctx context.Context, flat entity.Flat) (entity.Flat, error)
		Update(ctx context.Context, flat entity.Flat) (entity.Flat, error)
	}

	FlatRepo interface {
		Store(ctx context.Context, flat entity.Flat) (entity.Flat, error)
		Update(ctx context.Context, flat entity.Flat) (entity.Flat, error)
	}

	House interface {
		Create(ctx context.Context, house entity.House) (entity.House, error)
		Flats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error)
		//Subscribe()
	}

	HouseRepo interface {
		Store(ctx context.Context, house entity.House) (entity.House, error)
		GetFlats(ctx context.Context, houseID int, userType string) ([]entity.Flat, error)
	}

	Authorization interface {
		DummyLogin(ctx context.Context, userType string) (entity.TokenResponse, error)
		Login(ctx context.Context, request entity.LoginRequest) (entity.TokenResponse, error)
		Register(ctx context.Context, request entity.RegisterRequest) (entity.RegisterResponse, error)
	}

	AuthorizationRepo interface {
		GetUserByID(ctx context.Context, id string) (entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (entity.User, error)
		SaveUser(ctx context.Context, user entity.User) (entity.User, error)
	}
)
