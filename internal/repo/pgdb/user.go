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

func (r *UserRepo) UpdateUserPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	sql, args, _ := r.Builder.
		Update("users").
		Set("password", newPassword).
		Set("updated_at", "NOW()").
		Where("id = ? AND password = ?", userID, oldPassword).
		ToSql()

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.UpdateUserPassword - r.Pool.Exec: %v", err)
		return fmt.Errorf("UserRepo.UpdateUserPassword - r.Pool.Exec: %v", err)
	}

	if res.RowsAffected() == 0 {
		return repoerrs.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) UpdateUserByID(ctx context.Context, userID uuid.UUID, name, email, description *string, role *entity.RoleType) error {
	sqlBuilder := r.Builder.
		Update("users").
		Set("updated_at", "NOW()")

	if name != nil {
		sqlBuilder = sqlBuilder.Set("name", *name)
	}

	if email != nil {
		sqlBuilder = sqlBuilder.Set("email", *email)
	}

	if description != nil {
		sqlBuilder = sqlBuilder.Set("description", *description)
	}

	if role != nil {
		sqlBuilder = sqlBuilder.Set("role", *role)
	}

	sql, args, _ := sqlBuilder.
		Where("id = ?", userID).
		ToSql()

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.UpdateUserByID - r.Pool.Exec: %v", err)
		return fmt.Errorf("UserRepo.UpdateUserByID - r.Pool.Exec: %v", err)
	}

	if res.RowsAffected() == 0 {
		return repoerrs.ErrUserNotFound
	}

	return nil
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

func (r *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("userID = ?", userID).
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
		log.Errorf("UserRepo.GetUserByID - r.Pool.QueryRow: %v", err)
		if err == pgx.ErrNoRows {
			return entity.User{}, repoerrs.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByID - r.Pool.QueryRow: %v", err)
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

func (r *UserRepo) SetUserFollower(ctx context.Context, followerID uuid.UUID, followingID uuid.UUID) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("UserRepo.SetUserFollower - r.Pool.BeginTx: %v", err)
		return fmt.Errorf("UserRepo.SetUserFollower - r.Pool.BeginTx: %v", err)
	}

	sql, args, _ := r.Builder.
		Insert("users_followers").
		Columns("follower_id", "following_id").
		Values(followerID, followingID).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
		return fmt.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Update("users").
		Set("followers_count", "followers_count + 1").
		Where("id = ?", followingID).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
		return fmt.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Update("users").
		Set("following_count", "following_count + 1").
		Where("id = ?", followerID).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
		return fmt.Errorf("UserRepo.SetUserFollower - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("UserRepo.SetUserFollower - tx.Commit: %v", err)
		return fmt.Errorf("UserRepo.SetUserFollower - tx.Commit: %v", err)
	}

	return nil
}

func (r *UserRepo) GetUserFollowers(ctx context.Context, userID uuid.UUID) ([]entity.User, error) {
	sql, args, _ := r.Builder.
		Select("u.*").
		From("users_followers uf").
		Join("users u ON u.id = uf.follower_id").
		Where("uf.following_id = ?", userID).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.GetUserFollowers - r.Pool.Query: %v", err)
		return nil, fmt.Errorf("UserRepo.GetUserFollowers - r.Pool.Query: %v", err)
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
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
			log.Errorf("UserRepo.GetUserFollowers - rows.Scan: %v", err)
			return nil, fmt.Errorf("UserRepo.GetUserFollowers - rows.Scan: %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) GetUserFollowings(ctx context.Context, userID uuid.UUID) ([]entity.User, error) {
	sql, args, _ := r.Builder.
		Select("u.*").
		From("users_followers uf").
		Join("users u ON u.id = uf.following_id").
		Where("uf.follower_id = ?", userID).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("UserRepo.GetUserFollowings - r.Pool.Query: %v", err)
		return nil, fmt.Errorf("UserRepo.GetUserFollowings - r.Pool.Query: %v", err)
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
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
			log.Errorf("UserRepo.GetUserFollowings - rows.Scan: %v", err)
			return nil, fmt.Errorf("UserRepo.GetUserFollowings - rows.Scan: %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}
