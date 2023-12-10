package v1

import (
	"blog-backend/internal/entity"
	"blog-backend/internal/usecase"
	"github.com/labstack/echo/v4"
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
		cookie, err := c.Cookie("access-token")
		if err != nil {
			return echo.ErrForbidden
		}

		token := cookie.Value

		userID, role, err := h.authUseCase.ParseToken(c.Request().Context(), usecase.AuthParseTokenInput{
			Token: token,
		})
		if err != nil {
			return echo.ErrForbidden
		}

		c.Set(userIDCtx, userID)
		c.Set(userRoleCtx, role)

		return next(c)
	}
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
