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
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, password string) error
	UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	SetUserFollower(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error
	GetUserFollowers(ctx context.Context, userID uuid.UUID) ([]entity.User, error)
	GetUserFollowings(ctx context.Context, userID uuid.UUID) ([]entity.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
	}
}
