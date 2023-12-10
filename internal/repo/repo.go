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

type Article interface {
	CreateArticle(ctx context.Context, article entity.Article) (uuid.UUID, error)
	GetArticleByID(ctx context.Context, id uuid.UUID) (entity.Article, error)
	GetArticlesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]entity.Article, error)
	GetNewestArticles(ctx context.Context, limit, offset int) ([]entity.Article, error)
	SetArticleFavorite(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	RemoveArticleFavorite(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	GetFavoriteArticles(ctx context.Context, userID uuid.UUID) ([]entity.Article, error)
}

type Repositories struct {
	User
	Article
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepo(pg),
		Article: pgdb.NewArticleRepo(pg),
	}
}
