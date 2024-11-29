package application

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net"

	"github.com/SpaceSlow/test-task-backend-junior-medods/generated/openapi"
	usersapp "github.com/SpaceSlow/test-task-backend-junior-medods/internal/application/users"
	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type UserService interface {
	Tokens(userGUID uuid.UUID, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
	RefreshTokens(access *users.AccessToken, refresh *users.RefreshToken, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
}

type NotifierService interface {
	SendSuspiciousActivityMail(email string, newIP net.IP) error
}

type httpServer struct {
	usersapp.UserHandlers
}

func SetupHTTPServer(userService UserService) *echo.Echo {
	server := &httpServer{}
	server.UserHandlers = usersapp.SetupHandlers(userService)

	router := echo.New()
	api := router.Group("/api")
	openapi.RegisterHandlers(api, server)

	return router
}
