package service

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/pkg/hasher"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
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

var (
	ErrCannotCreateUser = fmt.Errorf("cannot create user")
	ErrCannotGetUser    = fmt.Errorf("cannot get user")
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")
	ErrTokenClaimsType  = fmt.Errorf("token claims are not of type TokenClaims")
)

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

	userId, err := s.authRepo.CreateUser(ctx, user)
	if err != nil {
		log.Errorf("AuthService.CreateUser - c.authRepo.CreateUser: %v", err)
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	// get user from DB
	user, err := s.authRepo.GetUser(ctx, input.Username, s.passwordHasher.Hash(input.Password))
	if err != nil {
		log.Errorf("AuthService.GenerateToken: cannot get user: %v", err)
		return "", ErrCannotGetUser
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
		log.Errorf("AuthService.GenerateToken: cannot sign token: %v", err)
		return "", ErrCannotSignToken
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
		log.Errorf("AuthService.ParseToken: cannot parse token: %v", err)
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrTokenClaimsType
	}

	return claims.UserId, nil
}
