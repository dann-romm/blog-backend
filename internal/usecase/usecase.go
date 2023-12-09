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
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (uuid.UUID, entity.RoleType, error)
}

type User interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (uuid.UUID, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, password string) error
	UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error
}

type UseCases struct {
	Auth Auth
	User User
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
		User: NewUserUseCase(deps.Repos, deps.Hasher),
	}
}
