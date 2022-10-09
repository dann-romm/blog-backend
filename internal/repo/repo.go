package repo

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo/pgdb"
	"blog-backend/pkg/postgres"
	"context"
)

type Auth interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUser(ctx context.Context, username, passwordHash string) (entity.User, error)
}

type Repositories struct {
	Auth
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Auth: pgdb.NewAuthRepo(pg),
	}
}
