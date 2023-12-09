package usecase

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/repo"
	"blog-backend/internal/repo/repoerrs"
	"blog-backend/pkg/hasher"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
}

var (
	ErrUserAlreadyExists    = fmt.Errorf("user already exists")
	ErrCannotCreateUser     = fmt.Errorf("cannot create user")
	ErrCannotChangePassword = fmt.Errorf("cannot change password")
	ErrCannotChangeEmail    = fmt.Errorf("cannot change email")
)

func NewUserUseCase(userRepo repo.User, passwordHasher hasher.PasswordHasher) User {
	return &UserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, input AuthCreateUserInput) (uuid.UUID, error) {
	user := entity.User{
		Name:     input.Name,
		Username: input.Username,
		Password: u.passwordHasher.Hash(input.Password),
		Email:    input.Email,
		Role:     entity.RoleUser,
	}

	userID, err := u.userRepo.CreateUser(ctx, user)
	if err == repoerrs.ErrUserAlreadyExists {
		return uuid.UUID{}, ErrUserAlreadyExists
	}
	if err != nil {
		return uuid.UUID{}, ErrCannotCreateUser
	}
	return userID, nil
}

func (u *UserUseCase) UpdateUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	err := u.userRepo.UpdateUserPassword(ctx, userID, u.passwordHasher.Hash(password))
	if err != nil {
		return ErrCannotChangePassword
	}
	return nil
}

func (u *UserUseCase) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error {
	err := u.userRepo.UpdateUserEmail(ctx, userID, email)
	if err != nil {
		return ErrCannotChangeEmail
	}
	return nil
}
