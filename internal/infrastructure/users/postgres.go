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
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &PostgresRepo{
		pool: pool,
		ctx:  ctx,
	}, nil
}

func (r *PostgresRepo) CreateRefreshToken(userGUID uuid.UUID, hash string) error {
	_, err := r.pool.Exec(r.ctx, "UPDATE users SET refresh_token=$1 WHERE id=$2", hash, userGUID)
	return err
}

func (r *PostgresRepo) RefreshToken(userGUID uuid.UUID) (*users.RefreshToken, error) {
	row := r.pool.QueryRow(r.ctx, "SELECT id, refresh_token FROM users WHERE id=$1", userGUID)
	var (
		id           uuid.UUID
		refreshToken string
	)
	err := row.Scan(&id, &refreshToken)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, users.NewNoUserError(userGUID)
	} else if err != nil {
		return nil, err
	}
	refresh := users.RefreshToken(refreshToken)

	return &refresh, nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}
