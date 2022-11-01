package pgdb

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo/repoerrs"
	"blog-backend/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("name", "username", "password", "email").
		Values(user.Name, user.Username, user.Password, user.Email).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return uuid.UUID{}, repoerrs.ErrUserAlreadyExists
			}
		}
		return uuid.UUID{}, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("username = ? AND password = ?", username, password).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
		&user.Description,
		&user.ArticlesCount,
		&user.CommentsCount,
		&user.FavoritesArticlesCount,
		&user.FavoritesCommentsCount,
		&user.FollowersCount,
		&user.FollowingCount,
	)
	if err != nil {
		log.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
		if err == pgx.ErrNoRows {
			return entity.User{}, repoerrs.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("id = ?", id).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
		&user.Description,
		&user.ArticlesCount,
		&user.CommentsCount,
		&user.FavoritesArticlesCount,
		&user.FavoritesCommentsCount,
		&user.FollowersCount,
		&user.FollowingCount,
	)
	if err != nil {
		log.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
		if err == pgx.ErrNoRows {
			return entity.User{}, repoerrs.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("username = ?", username).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
		&user.Description,
		&user.ArticlesCount,
		&user.CommentsCount,
		&user.FavoritesArticlesCount,
		&user.FavoritesCommentsCount,
		&user.FollowersCount,
		&user.FollowingCount,
	)
	if err != nil {
		log.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
		if err == pgx.ErrNoRows {
			return entity.User{}, repoerrs.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
