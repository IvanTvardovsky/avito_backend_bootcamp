package usecase

import (
	"avito_bootcamp/internal/entity"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthUseCase struct {
	repo         AuthorizationRepo
	jwtSecretKey string
}

var _ Authorization = (*AuthUseCase)(nil)

func NewAuthUseCase(r AuthorizationRepo, jwtSecretKey string) *AuthUseCase {
	return &AuthUseCase{
		repo:         r,
		jwtSecretKey: jwtSecretKey,
	}
}

func (a *AuthUseCase) DummyLogin(ctx context.Context, userType string) (entity.TokenResponse, error) {
	if userType != "client" && userType != "moderator" {
		return entity.TokenResponse{}, fmt.Errorf("invalid user type")
	}

	tokenString, err := a.createToken("", userType)
	if err != nil {
		return entity.TokenResponse{}, err
	}

	return entity.TokenResponse{Token: tokenString}, nil
}

func (a *AuthUseCase) Login(ctx context.Context, request entity.LoginRequest) (entity.TokenResponse, error) {
	user, err := a.repo.GetUserByID(ctx, request.ID)
	if err != nil {
		return entity.TokenResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return entity.TokenResponse{}, fmt.Errorf("invalid credentials")
	}

	tokenString, err := a.createToken(user.ID, user.UserType)
	if err != nil {
		return entity.TokenResponse{}, err
	}

	return entity.TokenResponse{Token: tokenString}, nil
}

func (a *AuthUseCase) Register(ctx context.Context, request entity.RegisterRequest) (entity.RegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.RegisterResponse{}, err
	}

	user := entity.User{
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		UserType:     request.UserType,
		CreatedAt:    time.Now(),
	}

	user, err = a.repo.SaveUser(ctx, user)
	if err != nil {
		return entity.RegisterResponse{}, err
	}

	return entity.RegisterResponse{UserID: user.ID}, nil
}

func (a *AuthUseCase) createToken(userID, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}
	if userID != "" {
		claims["user_id"] = userID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecretKey))
}
