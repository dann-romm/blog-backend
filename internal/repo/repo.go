package repo

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo/pgdb"
	"blog-backend/pkg/postgres"
	"context"
	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
	}
}
