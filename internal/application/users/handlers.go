package users

import (
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/SpaceSlow/test-task-backend-junior-medods/generated/openapi"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type UserService interface {
	Tokens(userGUID uuid.UUID, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
	RefreshTokens(access *users.AccessToken, refresh *users.RefreshToken, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
}

type UserHandlers struct {
	service UserService
}

func SetupHandlers(userService UserService) UserHandlers {
	return UserHandlers{service: userService}
}

func (h UserHandlers) PostUsersRefresh(c echo.Context) error {
	var req openapi.TokensRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Error: err.Error(),
		})
	}
	accessToken := users.AccessToken(req.Access)
	refreshToken, err := users.ParseRefreshToken(req.Refresh)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Error: err.Error(),
		})
	}

	access, refresh, err := h.service.RefreshTokens(
		&accessToken,
		refreshToken,
		net.ParseIP(c.RealIP()),
	)
	if errors.Is(errors.Join(users.ErrNoRefreshToken, users.ErrInvalidAccessToken, users.ErrInvalidRefreshToken), err) {
		return c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Error: err.Error(),
		})
	} else if err != nil {
		slog.Error("internal error", slog.String("error", err.Error()))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, openapi.TokensResponse{
		Access:  access.String(),
		Refresh: refresh.String(),
	})
}

func (h UserHandlers) GetUsersTokens(ctx echo.Context, params openapi.GetUsersTokensParams) error {
	access, refresh, err := h.service.Tokens(params.Guid, net.ParseIP(ctx.RealIP()))

	var errNoUser users.NoUserError
	if errors.As(err, &errNoUser) {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Error: errNoUser.Error(),
		})
	} else if err != nil {
		slog.Error("internal error", slog.String("error", err.Error()))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, openapi.TokensResponse{
		Access:  access.String(),
		Refresh: refresh.String(),
	})
}
