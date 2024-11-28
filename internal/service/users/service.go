package users

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type Repository interface {
	CreateRefreshToken(userGUID uuid.UUID, refresh *users.RefreshToken) error
	EmailByUUID(userGUID uuid.UUID) (string, error)
	UserByEmail(email string) (*users.User, error)
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
	const method = "UserService.Tokens"

	email, err := s.repo.EmailByUUID(userGUID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, users.NewNoUserError(userGUID)
	} else if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", method, err)
	}

	refresh, err := users.NewRefreshToken()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", method, err)
	}
	err = s.repo.CreateRefreshToken(userGUID, refresh)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", method, err)
	}

	access, err := users.NewAccessToken(email, ip, s.cfg.TokenLifetime(), s.cfg.SecretKey())
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", method, err)
	}

	return access, refresh, nil
}

func (s *UserService) RefreshTokens(access *users.AccessToken, refresh *users.RefreshToken, ip net.IP) (*users.AccessToken, *users.RefreshToken, error) {
	const method = "UserService.RefreshTokens"

	accessTokenIP, err := access.ParseIP(s.cfg.SecretKey())
	if errors.Is(err, users.ErrInvalidAccessToken) {
		return nil, nil, users.ErrInvalidAccessToken
	}

	if !accessTokenIP.Equal(ip) {
		// TODO: sending warning mail
	}
	accessTokenEmail, err := access.ParseEmail(s.cfg.SecretKey())

	user, err := s.repo.UserByEmail(accessTokenEmail)
	if errors.Is(err, users.ErrNoRefreshToken) {
		return nil, nil, users.ErrNoRefreshToken
	} else if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", method, err)
	}
	if !refresh.Valid(user.RefreshTokenHash()) {
		return nil, nil, users.ErrInvalidRefreshToken
	}

	return s.Tokens(user.Id(), ip)
}
