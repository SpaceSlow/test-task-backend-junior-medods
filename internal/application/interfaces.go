package application

import (
	"net"

	"github.com/google/uuid"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type UserService interface {
	Tokens(userGUID uuid.UUID, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
	RefreshTokens(access *users.AccessToken, refresh *users.RefreshToken, ip net.IP) (*users.AccessToken, *users.RefreshToken, error)
}

type NotifierService interface {
	SendSuspiciousActivityMail(email string, newIP net.IP) error
}
