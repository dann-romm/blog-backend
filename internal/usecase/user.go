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
	ErrUserAlreadyExists               = fmt.Errorf("user already exists")
	ErrCannotCreateUser                = fmt.Errorf("cannot create user")
	ErrHaveNoPermission                = fmt.Errorf("have no permission")
	ErrCannotUpdatePasswordToIdentical = fmt.Errorf("cannot update password to identical")
	ErrNothingToUpdate                 = fmt.Errorf("nothing to update")
)

func NewUserUseCase(userRepo repo.User, passwordHasher hasher.PasswordHasher) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, input UserCreateUserInput) (uuid.UUID, error) {
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

func (u *UserUseCase) UpdateUser(ctx context.Context, input UserUpdateUserInput) error {
	if input.Name == nil && input.Email == nil && input.Role == nil && input.Description == nil {
		return ErrNothingToUpdate
	}

	err := u.checkPermissions(ctx, input)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdateUser(
		ctx,
		input.UserID,
		input.Name,
		input.Email,
		input.Role,
		input.Description,
	)
	if err == repoerrs.ErrUserNotFound {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) UpdateUserPassword(ctx context.Context, input UserUpdateUserPasswordInput) error {
	if input.NewPassword == input.OldPassword {
		return ErrCannotUpdatePasswordToIdentical
	}

	err := u.userRepo.UpdateUserPassword(
		ctx,
		input.UserID,
		u.passwordHasher.Hash(input.OldPassword),
		u.passwordHasher.Hash(input.NewPassword),
	)
	if err == repoerrs.ErrUserNotFound {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) checkPermissions(ctx context.Context, input UserUpdateUserInput) error {
	// anyone can update himself except role
	if input.UserID == input.RequestedUserID {
		if input.Role != nil {
			return ErrHaveNoPermission
		}
		return nil
	}

	// user can't update other users
	if input.RequestedUserRole == entity.RoleUser {
		return ErrHaveNoPermission
	}

	// moderator can't update other emails
	if input.RequestedUserRole == entity.RoleModerator && input.Email != nil {
		return ErrHaveNoPermission
	}

	user, err := u.userRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if input.RequestedUserRole == entity.RoleModerator {
		// moderator can't update any roles
		if input.Role != nil {
			return ErrHaveNoPermission
		}

		// moderator can't update other moderators and admins
		if user.Role == entity.RoleModerator || user.Role == entity.RoleAdmin {
			return ErrHaveNoPermission
		}
	}

	if input.RequestedUserRole == entity.RoleAdmin {
		// admin can't update other admins
		if user.Role == entity.RoleAdmin {
			return ErrHaveNoPermission
		}

		// admin can't update role to admin
		if input.Role != nil && entity.RoleType(*input.Role) == entity.RoleAdmin {
			return ErrHaveNoPermission
		}
	}

	return nil
}
