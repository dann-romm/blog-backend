package pgdb

import (
	"blog-backend/internal/entity"
	"blog-backend/pkg/postgres"
	"context"
	"github.com/google/uuid"
)

type ArticleRepo struct {
	*postgres.Postgres
}

func NewArticleRepo(pg *postgres.Postgres) *ArticleRepo {
	return &ArticleRepo{pg}
}

func (a ArticleRepo) CreateArticle(ctx context.Context, article entity.Article) (uuid.UUID, error) {
	sql, args, _ := a.Builder.
		Insert("articles").
		Columns("author_id", "title", "description", "content").
		Values(article.AuthorID, article.Title, article.Description, article.Content).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (a ArticleRepo) GetArticleByID(ctx context.Context, id uuid.UUID) (entity.Article, error) {
	sql, args, _ := a.Builder.
		Select("*").
		From("articles").
		Where("id = ?", id).
		ToSql()

	var article entity.Article
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(
		&article.Id,
		&article.AuthorID,
		&article.Title,
		&article.Description,
		&article.Content,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.ViewsCount,
		&article.CommentsCount,
		&article.FavoritesCount,
		&article.VotesUpCount,
		&article.VotesDownCount,
	)
	if err != nil {
		return entity.Article{}, err
	}

	return article, nil
}

func (a ArticleRepo) GetArticlesByAuthorID(ctx context.Context, authorID uuid.UUID) ([]entity.Article, error) {
	sql, args, _ := a.Builder.
		Select("*").
		From("articles").
		Where("author_id = ?", authorID).
		ToSql()

	rows, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.Article
	for rows.Next() {
		var article entity.Article
		err := rows.Scan(
			&article.Id,
			&article.AuthorID,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.ViewsCount,
			&article.CommentsCount,
			&article.FavoritesCount,
			&article.VotesUpCount,
			&article.VotesDownCount,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (a ArticleRepo) GetNewestArticles(ctx context.Context, limit, offset int) ([]entity.Article, error) {
	sql, args, _ := a.Builder.
		Select("*").
		From("articles").
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.Article
	for rows.Next() {
		var article entity.Article
		err := rows.Scan(
			&article.Id,
			&article.AuthorID,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.ViewsCount,
			&article.CommentsCount,
			&article.FavoritesCount,
			&article.VotesUpCount,
			&article.VotesDownCount,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (a ArticleRepo) SetArticleFavorite(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	sql, args, _ := a.Builder.
		Insert("users_articles_favorites").
		Columns("user_id", "article_id").
		Values(userID, articleID).
		ToSql()

	_, err := a.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (a ArticleRepo) RemoveArticleFavorite(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	sql, args, _ := a.Builder.
		Delete("users_articles_favorites").
		Where("user_id = ?", userID).
		Where("article_id = ?", articleID).
		ToSql()

	_, err := a.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (a ArticleRepo) GetFavoriteArticles(ctx context.Context, userID uuid.UUID) ([]entity.Article, error) {
	sql, args, _ := a.Builder.
		Select("a.*").
		From("users_articles_favorites uf").
		Join("articles a ON a.id = uf.article_id").
		Where("uf.user_id = ?", userID).
		ToSql()

	rows, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []entity.Article
	for rows.Next() {
		var article entity.Article
		err := rows.Scan(
			&article.Id,
			&article.AuthorID,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.ViewsCount,
			&article.CommentsCount,
			&article.FavoritesCount,
			&article.VotesUpCount,
			&article.VotesDownCount,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}
