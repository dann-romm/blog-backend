package usecase

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/internal/repo/repoerrs"
	"blog-backend/pkg/hasher"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID       `json:"user_id"`
	Role   entity.RoleType `json:"role"`
}

type AuthUseCase struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

var (
	ErrCannotGetUser    = fmt.Errorf("cannot get user")
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")
	ErrTokenClaimsType  = fmt.Errorf("token claims are not of type TokenClaims")
	ErrUserNotFound     = fmt.Errorf("user not found")
)

func NewAuthUseCase(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (u *AuthUseCase) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	// get user from DB
	user, err := u.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, u.passwordHasher.Hash(input.Password))
	if err == repoerrs.ErrUserNotFound {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", ErrCannotGetUser
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(u.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
		Role:   user.Role,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(u.signKey))
	if err != nil {
		log.Errorf("AuthUseCase.GenerateToken: cannot sign token: %v", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (u *AuthUseCase) ParseToken(accessToken string) (uuid.UUID, entity.RoleType, error) {
	claims, err := u.parseToken(accessToken)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	return claims.UserID, claims.Role, nil
}

func (u *AuthUseCase) parseToken(accessToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(u.signKey), nil
	})

	if err != nil {
		return nil, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrTokenClaimsType
	}

	return claims, nil
}
