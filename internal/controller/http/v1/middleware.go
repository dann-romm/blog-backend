package v1

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	userIDCtx   = "userID"
	userRoleCtx = "userRole"
)

type AuthMiddleware struct {
	authUseCase usecase.Auth
}

func NewAuthMiddleware(authUseCase usecase.Auth) *AuthMiddleware {
	return &AuthMiddleware{authUseCase: authUseCase}
}

// Authorize - проверка авторизации пользователя
// если пользователь авторизован, то в контекст запроса добавляется его id и роль (user, moderator, admin)
func (h *AuthMiddleware) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized")
		}

		userID, role, err := h.authUseCase.ParseToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized")
		}

		c.Set(userIDCtx, userID)
		c.Set(userRoleCtx, role)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get(userRoleCtx).(entity.RoleType)
		if role != entity.RoleAdmin {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func ModeratorOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get(userRoleCtx).(entity.RoleType)
		if role != entity.RoleAdmin && role != entity.RoleModerator {
			return echo.ErrForbidden
		}

		return next(c)
	}
}

func UserOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get(userRoleCtx).(entity.RoleType)
		if role != entity.RoleAdmin && role != entity.RoleModerator && role != entity.RoleUser {
			return echo.ErrForbidden
		}

		return next(c)
	}
}
