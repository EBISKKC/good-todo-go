package presenter

import (
	"net/http"
	"time"

	"good-todo-go/internal/presentation/public/api"
	"good-todo-go/internal/usecase/output"

	"github.com/labstack/echo/v4"
)

type IAuthPresenter interface {
	Register(ctx echo.Context, out *output.AuthOutput) error
	Login(ctx echo.Context, out *output.AuthOutput) error
	VerifyEmail(ctx echo.Context, out *output.VerifyEmailOutput) error
	RefreshToken(ctx echo.Context, out *output.AuthOutput) error
}

type AuthPresenter struct{}

func NewAuthPresenter() IAuthPresenter {
	return &AuthPresenter{}
}

func (p *AuthPresenter) Register(ctx echo.Context, out *output.AuthOutput) error {
	return ctx.JSON(http.StatusCreated, toAuthResponse(out))
}

func (p *AuthPresenter) Login(ctx echo.Context, out *output.AuthOutput) error {
	return ctx.JSON(http.StatusOK, toAuthResponse(out))
}

func (p *AuthPresenter) VerifyEmail(ctx echo.Context, out *output.VerifyEmailOutput) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": out.Message,
	})
}

func (p *AuthPresenter) RefreshToken(ctx echo.Context, out *output.AuthOutput) error {
	return ctx.JSON(http.StatusOK, toAuthResponse(out))
}

func toAuthResponse(out *output.AuthOutput) *api.AuthResponse {
	role := api.UserResponseRole(out.User.Role)
	createdAt, _ := time.Parse(time.RFC3339, out.User.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, out.User.UpdatedAt)
	return &api.AuthResponse{
		AccessToken:  &out.AccessToken,
		RefreshToken: &out.RefreshToken,
		TokenType:    &out.TokenType,
		ExpiresIn:    &out.ExpiresIn,
		User: &api.UserResponse{
			Id:            &out.User.ID,
			Email:         &out.User.Email,
			Name:          &out.User.Name,
			Role:          &role,
			EmailVerified: &out.User.EmailVerified,
			TenantId:      &out.User.TenantID,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
		},
	}
}
