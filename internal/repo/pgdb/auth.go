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

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("name", "username", "password", "email").
		Values(user.Name, user.Username, user.Password, user.Email).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Errorf("AuthRepo.CreateUser - r.Builder: %v", err)
		return 0, fmt.Errorf("AuthRepo.CreateUser - r.Builder: %v", err)
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
		log.Errorf("AuthRepo.CreateUser - r.Pool.QueryRow: %v", err)
		return 0, fmt.Errorf("AuthRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *AuthRepo) GetUser(ctx context.Context, username, passwordHash string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "username", "password", "email").
		From("users").
		Where("username = ? AND password = ?", username, passwordHash).
		ToSql()

	if err != nil {
		log.Errorf("AuthRepo.GetUser - r.Builder: %v", err)
		return entity.User{}, fmt.Errorf("AuthRepo.GetUser - r.Builder: %v", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		log.Errorf("AuthRepo.GetUser - r.Pool.QueryRow: %v", err)
		if err == pgx.ErrNoRows {
			return entity.User{}, repoerrs.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("AuthRepo.GetUser - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
