package pgdb

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo/repoerrs"
	"blog-backend/pkg/postgres"
	"context"
	"errors"
	"fmt"
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

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("name", "username", "password", "email").
		Values(user.Name, user.Username, user.Password, user.Email).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Errorf("UserRepo.CreateUser - r.Builder: %v", err)
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Builder: %v", err)
	}

	var id int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrUserAlreadyExists
			}
		}
		log.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where("username = ? AND password = ?", username, password).
		ToSql()

	if err != nil {
		log.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Builder: %v", err)
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Builder: %v", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
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

func (r *UserRepo) GetUserById(ctx context.Context, id int) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		log.Errorf("UserRepo.GetUserById - r.Builder: %v", err)
		return entity.User{}, fmt.Errorf("UserRepo.GetUserById - r.Builder: %v", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
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
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where("username = ?", username).
		ToSql()

	if err != nil {
		log.Errorf("UserRepo.GetUserByUsername - r.Builder: %v", err)
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsername - r.Builder: %v", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
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
