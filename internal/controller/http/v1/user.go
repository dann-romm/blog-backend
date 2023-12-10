package v1

import (
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

	g.POST("/user/change-password", r.changePassword)
	g.POST("/user/change-email", r.changeEmail)
}

type changePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

func (r *userRoutes) changePassword(c echo.Context) error {
	var input changePasswordInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err = r.userUseCase.UpdateUserPassword(c.Request().Context(), usecase.UserUpdateUserPasswordInput{
		UserID:      c.Get(userIDCtx).(uuid.UUID),
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	})

	if err == usecase.ErrCannotChangePassword {
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

type changeEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *userRoutes) changeEmail(c echo.Context) error {
	var input changeEmailInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err = r.userUseCase.UpdateUserEmail(c.Request().Context(), usecase.UserUpdateUserEmailInput{
		UserID: c.Get(userIDCtx).(uuid.UUID),
		Email:  input.Email,
	})

	if err == usecase.ErrCannotChangeEmail {
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
