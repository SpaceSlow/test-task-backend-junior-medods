package users

import (
	"errors"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/SpaceSlow/test-task-backend-junior-medods/generated/openapi"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type UserService interface {
	Tokens(userGUID uuid.UUID, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
}

type UserHandlers struct {
	service UserService
}

func SetupHandlers(userService UserService) UserHandlers {
	return UserHandlers{service: userService}
}

func (h UserHandlers) PostUsersRefresh(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h UserHandlers) GetUsersTokens(ctx echo.Context, params openapi.GetUsersTokensParams) error {
	access, refresh, err := h.service.Tokens(params.Guid, net.IP(ctx.RealIP()))

	var errNoUser users.NoUserError
	if errors.As(err, &errNoUser) {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Error: errNoUser.Error(),
		})
	} else if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, openapi.TokensResponse{
		Access:  access.String(),
		Refresh: refresh.String(),
	})
}
