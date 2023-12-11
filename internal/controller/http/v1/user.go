package v1

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userRoutes struct {
	userUseCase usecase.User
}

func newUserRoutes(g *echo.Group, userUseCase usecase.User) {
	r := &userRoutes{
		userUseCase: userUseCase,
	}

	g.PUT("/users", r.updateUser)
	g.PUT("/users/password", r.updateUserPassword)
}

type updateUserInput struct {
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	Name        *string   `json:"name" validate:"omitempty,min=3,max=256"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Role        *string   `json:"role"`
	Description *string   `json:"description"`
}

func (r *userRoutes) updateUser(c echo.Context) error {
	var input updateUserInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	requestedUserRole := c.Get(userRoleCtx).(entity.RoleType)
	requestedUserID := c.Get(userIDCtx).(uuid.UUID)

	err = r.userUseCase.UpdateUser(c.Request().Context(), usecase.UserUpdateUserInput{
		RequestedUserID:   requestedUserID,
		RequestedUserRole: requestedUserRole,
		UserID:            input.UserID,
		Name:              input.Name,
		Email:             input.Email,
		Role:              input.Role,
		Description:       input.Description,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}

type updateUserPasswordInput struct {
	OldPassword string `json:"old_password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

func (r *userRoutes) updateUserPassword(c echo.Context) error {
	var input updateUserPasswordInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	userID := c.Get(userIDCtx).(uuid.UUID)

	err = r.userUseCase.UpdateUserPassword(c.Request().Context(), usecase.UserUpdateUserPasswordInput{
		UserID:      userID,
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	})

	if err == usecase.ErrUserNotFound {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err == usecase.ErrCannotUpdatePasswordToIdentical {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
