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

func (u *UserUseCase) GetUserByUsername(ctx context.Context, input UserGetUserByUsernameInput) (entity.User, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, input.Username)
	if err == repoerrs.ErrUserNotFound {
		return entity.User{}, ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, input UserUpdateUserInput) error {
	if input.NewName == nil && input.NewEmail == nil && input.NewRole == nil && input.NewDescription == nil {
		return ErrNothingToUpdate
	}

	user, err := u.userRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return err
	}

	// check updating the same
	if (input.NewName == nil || *input.NewName == user.Name) &&
		(input.NewEmail == nil || *input.NewEmail == user.Email) &&
		(input.NewRole == nil || *input.NewRole == user.Role) &&
		(input.NewDescription == nil || *input.NewDescription == user.Description) {
		return ErrNothingToUpdate
	}

	err = u.checkPermissions(user, input.RequestedUserID, input.RequestedUserRole, input.NewRole, input.NewEmail)
	if err != nil {
		return err
	}

	err = u.userRepo.UpdateUserByID(
		ctx,
		user.ID,
		input.NewName,
		input.NewEmail,
		input.NewDescription,
		input.NewRole,
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

func (u *UserUseCase) checkPermissions(
	user entity.User,
	requestedUserID uuid.UUID,
	requestedUserRole entity.RoleType,
	newRole *entity.RoleType,
	newEmail *string,
) error {
	// anyone can update himself except role
	if user.ID == requestedUserID {
		if newRole != nil {
			return ErrHaveNoPermission
		}
		return nil
	}

	// user can't update other users
	if requestedUserRole == entity.RoleUser {
		return ErrHaveNoPermission
	}

	// moderator can't update other emails
	if requestedUserRole == entity.RoleModerator && newEmail != nil {
		return ErrHaveNoPermission
	}

	if requestedUserRole == entity.RoleModerator {
		// moderator can't update any roles
		if newRole != nil {
			return ErrHaveNoPermission
		}

		// moderator can't update other moderators and admins
		if user.Role == entity.RoleModerator || user.Role == entity.RoleAdmin {
			return ErrHaveNoPermission
		}
	}

	if requestedUserRole == entity.RoleAdmin {
		// admin can't update other admins
		if user.Role == entity.RoleAdmin {
			return ErrHaveNoPermission
		}

		// admin can't update role to admin
		if newRole != nil && *newRole == entity.RoleAdmin {
			return ErrHaveNoPermission
		}
	}

	return nil
}
