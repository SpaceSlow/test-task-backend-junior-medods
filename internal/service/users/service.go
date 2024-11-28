package users

import (
	"errors"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type Repository interface {
	CreateRefreshToken(userGUID uuid.UUID, refresh *users.RefreshToken) error
	RefreshToken(userGUID uuid.UUID) (*users.RefreshToken, error)
	UserGUID(*users.RefreshToken) (uuid.UUID, error)
}

type Config interface {
	TokenLifetime() time.Duration
	SecretKey() string
}

type UserService struct {
	repo Repository
	cfg  Config
}

func NewUserService(repo Repository, cfg Config) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *UserService) Tokens(userGUID uuid.UUID, ip net.IP) (*users.AccessToken, *users.RefreshToken, error) {
	refresh, err := s.repo.RefreshToken(userGUID)
	if errors.Is(err, users.ErrNoRefreshToken) {
		refresh, err = users.NewRefreshToken()
		if err != nil {
			return nil, nil, err
		}
		err = s.repo.CreateRefreshToken(userGUID, refresh)
	} else if err != nil {
		return nil, nil, err
	}

	access, err := users.NewAccessToken(ip, s.cfg.TokenLifetime(), s.cfg.SecretKey())
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (s *UserService) RefreshTokens(access *users.AccessToken, refresh *users.RefreshToken, ip net.IP) (*users.AccessToken, *users.RefreshToken, error) {
	accessTokenIP, err := access.IP(s.cfg.SecretKey())

	if !accessTokenIP.Equal(ip) {
		// TODO: sending warning mail
	}

	userGUID, err := s.repo.UserGUID(refresh)
	var errNoUser users.NoUserError
	if errors.As(err, &errNoUser) {
		return nil, nil, users.ErrInvalidRefreshToken
	} else if err != nil {
		return nil, nil, err
	}

	return s.Tokens(userGUID, ip)

}
