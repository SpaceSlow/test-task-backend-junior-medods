package users

import (
	"github.com/labstack/echo/v4"

	"github.com/SpaceSlow/test-task-backend-junior-medods/generated/openapi"
)

type UserHandlers struct {
}

func SetupHandlers() UserHandlers {
	return UserHandlers{}
}

func (h UserHandlers) PostUsersRefresh(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h UserHandlers) GetUsersTokens(ctx echo.Context, params openapi.GetUsersTokensParams) error {
	//TODO implement me
	panic("implement me")
}
