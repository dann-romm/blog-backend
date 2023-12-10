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

type Article interface {
	CreateArticle(ctx context.Context, input ArticleCreateArticleInput) (uuid.UUID, error)
	GetArticleByID(ctx context.Context, input ArticleGetArticleByIDInput) (entity.Article, error)
	GetArticlesByAuthorID(ctx context.Context, input ArticleGetArticlesByAuthorIDInput) ([]entity.Article, error)
	GetNewestArticles(ctx context.Context, input ArticleGetNewestArticlesInput) ([]entity.Article, error)
	SetArticleFavorite(ctx context.Context, input ArticleSetArticleFavoriteInput) error
	RemoveArticleFavorite(ctx context.Context, input ArticleRemoveArticleFavoriteInput) error
	GetFavoriteArticles(ctx context.Context, input ArticleGetFavoriteArticlesInput) ([]entity.Article, error)
}

type UseCases struct {
	Auth    Auth
	User    User
	Article Article
}

type UseCasesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewUseCases(deps UseCasesDependencies) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(deps.Repos, deps.Hasher, deps.SignKey, deps.TokenTTL),
		User:    NewUserUseCase(deps.Repos, deps.Hasher),
		Article: NewArticleUseCase(deps.Repos),
	}
}
