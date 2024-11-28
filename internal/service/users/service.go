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
