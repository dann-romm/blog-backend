package usecase

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ArticleUseCase struct {
	articleRepo repo.Article
}

var (
	ErrCannotCreateArticle = fmt.Errorf("cannot create article")
)

func NewArticleUseCase(articleRepo repo.Article) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (a *ArticleUseCase) CreateArticle(ctx context.Context, input ArticleCreateArticleInput) (uuid.UUID, error) {
	article := entity.Article{
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Description: input.Description,
		Content:     input.Content,
	}

	articleID, err := a.articleRepo.CreateArticle(ctx, article)
	if err != nil {
		return uuid.UUID{}, ErrCannotCreateArticle
	}
	return articleID, nil
}

func (a *ArticleUseCase) GetArticleByID(ctx context.Context, input ArticleGetArticleByIDInput) (entity.Article, error) {
	article, err := a.articleRepo.GetArticleByID(ctx, input.ID)
	if err != nil {
		return entity.Article{}, err
	}
	return article, nil
}

func (a *ArticleUseCase) GetArticlesByAuthorID(ctx context.Context, input ArticleGetArticlesByAuthorIDInput) ([]entity.Article, error) {
	articles, err := a.articleRepo.GetArticlesByAuthorID(ctx, input.AuthorID)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleUseCase) GetNewestArticles(ctx context.Context, input ArticleGetNewestArticlesInput) ([]entity.Article, error) {
	articles, err := a.articleRepo.GetNewestArticles(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleUseCase) SetArticleFavorite(ctx context.Context, input ArticleSetArticleFavoriteInput) error {
	err := a.articleRepo.SetArticleFavorite(ctx, input.UserID, input.ArticleID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleUseCase) RemoveArticleFavorite(ctx context.Context, input ArticleRemoveArticleFavoriteInput) error {
	err := a.articleRepo.RemoveArticleFavorite(ctx, input.UserID, input.ArticleID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleUseCase) GetFavoriteArticles(ctx context.Context, input ArticleGetFavoriteArticlesInput) ([]entity.Article, error) {
	articles, err := a.articleRepo.GetFavoriteArticles(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	return articles, nil
}
