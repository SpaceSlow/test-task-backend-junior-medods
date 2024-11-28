package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"
)

type PostgresRepo struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, dsn string) (*PostgresRepo, error) {
	const method = "NewPostgresRepo"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	return &PostgresRepo{
		pool: pool,
		ctx:  ctx,
	}, nil
}

func (r *PostgresRepo) CreateRefreshToken(userGUID uuid.UUID, refresh *users.RefreshToken) error {
	const method = "PostgresRepo.CreateRefreshToken"

	hash, err := refresh.Hash()
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}
	_, err = r.pool.Exec(r.ctx, "UPDATE users SET refresh_token=$1 WHERE id=$2", hash, userGUID)
	return err
}

func (r *PostgresRepo) RefreshToken(userGUID uuid.UUID) (*users.RefreshToken, error) {
	const method = "PostgresRepo.RefreshToken"

	row := r.pool.QueryRow(r.ctx, "SELECT refresh_token FROM users WHERE id=$1", userGUID)
	var refreshToken string
	err := row.Scan(&refreshToken)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, users.NewNoUserError(userGUID)
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	refresh := users.RefreshToken(refreshToken)

	return &refresh, nil
}

func (r *PostgresRepo) UserGUID(refresh *users.RefreshToken) (uuid.UUID, error) {
	const method = "PostgresRepo.UserGUID"

	row := r.pool.QueryRow(r.ctx, "SELECT id FROM users WHERE refresh_token=$1", refresh.String())
	var userGUID uuid.UUID
	err := row.Scan(&userGUID)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, users.ErrNoRefreshToken
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", method, err)
	}

	return userGUID, nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}
