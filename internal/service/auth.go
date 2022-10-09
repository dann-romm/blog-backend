package service

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/pkg/hasher"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	authRepo       repo.Auth
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(authRepo repo.Auth, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		authRepo:       authRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error) {
	user := entity.User{
		Name:     input.Name,
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
		Email:    input.Email,
	}

	return s.authRepo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	// get user from DB
	user, err := s.authRepo.GetUser(ctx, username, s.passwordHasher.Hash(password))
	if err != nil {
		return "", err
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, fmt.Errorf("token claims are not of type TokenClaims")
	}

	return claims.UserId, nil
}
