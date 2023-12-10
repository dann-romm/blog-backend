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

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type AuthParseTokenInput struct {
	Token string
}

type UserCreateUserInput struct {
	Name     string
	Username string
	Password string
	Email    string
}

type UserUpdateUserPasswordInput struct {
	UserID   uuid.UUID
	Password string
}

type UserUpdateUserEmailInput struct {
	UserID uuid.UUID
	Email  string
}

type Auth interface {
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(ctx context.Context, input AuthParseTokenInput) (uuid.UUID, entity.RoleType, error)
	GetTokenTTL() (time.Duration, error)
}

type User interface {
	CreateUser(ctx context.Context, input UserCreateUserInput) (uuid.UUID, error)
	UpdateUserPassword(ctx context.Context, input UserUpdateUserPasswordInput) error
	UpdateUserEmail(ctx context.Context, input UserUpdateUserEmailInput) error
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
