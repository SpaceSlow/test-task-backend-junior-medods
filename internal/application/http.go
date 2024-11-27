package application

import (
	"github.com/labstack/echo/v4"

	"github.com/SpaceSlow/test-task-backend-junior-medods/generated/openapi"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/application/users"
)

type httpServer struct {
	users.UserHandlers
}

func SetupHTTPServer(userService UserService) *echo.Echo {
	server := &httpServer{}
	server.UserHandlers = users.SetupHandlers(userService)

	router := echo.New()
	api := router.Group("/api")
	openapi.RegisterHandlers(api, server)

	return router
}
