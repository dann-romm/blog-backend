package v1

import (
	"blog-backend/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
	Email    string `json:"email" validate:"required,email"`
}

// регистрация пользователя
func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		log.Errorf("authRoutes.signUp: c.Bind: %v", err)
		newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Errorf("authRoutes.signUp: c.Validate: %v", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err == service.ErrUserAlreadyExists {
		log.Errorf("authRoutes.signUp: r.authService.CreateUser: %v", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err != nil {
		log.Errorf("authRoutes.signUp: r.authService.CreateUser: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
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

	if err := c.Bind(&input); err != nil {
		log.Errorf("authRoutes.signIn: c.Bind: %v", err)
		newErrorResponse(c, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Errorf("authRoutes.signIn: c.Validate: %v", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err == service.ErrUserNotFound {
		log.Errorf("authRoutes.signIn: r.authService.GenerateToken: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "invalid username or password")
		return err
	}
	if err != nil {
		log.Errorf("authRoutes.signIn: r.authService.GenerateToken: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
