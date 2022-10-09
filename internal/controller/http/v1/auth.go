package v1

import (
	"blog-backend/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type signUpInput struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Username string `json:"username" validate:"required,regexp=^[a-zA-Z0-9][a-zA-Z0-9_]{5,30}$"`
	Password string `json:"password" validate:"required,regexp=^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})"`
	Email    string `json:"email" validate:"required,email"`
}

// регистрация пользователя
func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// аутентификация пользователя
func (r *authRoutes) signIn(c echo.Context) error {
	// TODO: реализовать аутентификацию пользователя
	return nil
}
