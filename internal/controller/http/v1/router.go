package v1

import (
	"blog-backend/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func NewRouter(handler *echo.Echo, useCases *usecase.UseCases) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })
	handler.Static("/swagger-ui", "docs/swagger-ui")

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, useCases.Auth, useCases.User)
	}

	authMiddleware := NewAuthMiddleware(useCases.Auth)
	v1 := handler.Group("/api/v1", authMiddleware.Authorize)
	{
		newUserRoutes(v1, useCases.User)
		newArticleRoutes(v1, useCases.Article)
	}
}
