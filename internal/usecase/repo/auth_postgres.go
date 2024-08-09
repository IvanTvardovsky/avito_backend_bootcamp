package repo

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/pkg/postgres"
	"context"
	"github.com/google/uuid"
)

type AuthRepo struct {
	pg *postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{
		pg: pg,
	}
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password_hash, user_type, created_at FROM users WHERE id=$1`
	err := r.pg.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserType, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password_hash, user_type, created_at FROM users WHERE email=$1`
	err := r.pg.Pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.UserType, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthRepo) SaveUser(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, user_type, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pg.Pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.UserType, user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
