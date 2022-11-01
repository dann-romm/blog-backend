package v1

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	userIDCtx   = "userID"
	userRoleCtx = "userRole"
)

type AuthMiddleware struct {
	authService service.Auth
}

func NewAuthMiddleware(authService service.Auth) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// Authorize - проверка авторизации пользователя
// если пользователь авторизован, то в контекст запроса добавляется его id и роль (user, moderator, admin)
// если пользователь не авторизован, то в контекст запроса добавляется роль guest
func (h *AuthMiddleware) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			c.Set(userRoleCtx, entity.RoleGuest)
			return next(c)
		}

		userID, role, err := h.authService.ParseToken(token)
		if err != nil {
			c.Set(userRoleCtx, entity.RoleGuest)
			return next(c)
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
