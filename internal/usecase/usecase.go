package usecase

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/pkg/hasher"
	"context"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=usecase.go -destination=mocks/usecase.go -package=mocks

type AuthCreateUserInput struct {
	Name     string
	Username string
	Password string
	Email    string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (uuid.UUID, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (uuid.UUID, entity.RoleType, error)
}

type UseCases struct {
	Auth Auth
}

type UseCasesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewUseCases(deps UseCasesDependencies) *UseCases {
	return &UseCases{
		Auth: NewAuthUseCase(deps.Repos, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
