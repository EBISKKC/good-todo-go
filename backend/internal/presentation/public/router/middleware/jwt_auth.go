package middleware

import (
	"net/http"
	"strings"

	"good-todo-go/internal/pkg"
	"good-todo-go/internal/presentation/public/router/context_keys"

	"github.com/labstack/echo/v4"
)

func JWTAuth(jwtService *pkg.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := parts[1]
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			if claims.TokenType != pkg.AccessToken {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token type")
			}

			c.Set(context_keys.UserIDContextKey, claims.UserID)
			c.Set(context_keys.TenantIDContextKey, claims.TenantID)
			c.Set(context_keys.EmailContextKey, claims.Email)
			c.Set(context_keys.RoleContextKey, claims.Role)

			return next(c)
		}
	}
}
