package service

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/pkg/hasher"
	"context"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks

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

type Services struct {
	Auth Auth
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
