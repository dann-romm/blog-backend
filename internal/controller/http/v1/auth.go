package v1

import (
	"blog-backend/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authRoutes struct {
	authUseCase usecase.Auth
	userUseCase usecase.User
}

func newAuthRoutes(g *echo.Group, authUseCase usecase.Auth, userUseCase usecase.User) {
	r := &authRoutes{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type signUpInput struct {
	Name     string `json:"name" validate:"required,min=4,max=256"`
	Username string `json:"username" validate:"required,min=4,max=256"`
	Password string `json:"password" validate:"required,password"`
	Email    string `json:"email" validate:"required,email"`
}

// регистрация пользователя
func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.userUseCase.CreateUser(c.Request().Context(), usecase.UserCreateUserInput{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err == usecase.ErrUserAlreadyExists {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// аутентификация пользователя
func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput

	err := BindAndValidate(c, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authUseCase.GenerateToken(c.Request().Context(), usecase.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err == usecase.ErrUserNotFound {
		newErrorResponse(c, http.StatusBadRequest, "invalid username or password")
		return err
	}
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	ttl, _ := r.authUseCase.GetTokenTTL()
	c.SetCookie(&http.Cookie{
		Name:   "access-token",
		Value:  token,
		Path:   "/",
		MaxAge: int(ttl.Seconds()),
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token,
	})
}
