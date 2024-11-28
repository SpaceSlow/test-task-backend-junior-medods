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

	hash, err := refresh.GenerateHash()
	if err != nil {
		return fmt.Errorf("%s: %w", method, err)
	}
	_, err = r.pool.Exec(r.ctx, `UPDATE users SET refresh_token=$1 WHERE id=$2`, hash, userGUID)
	return err
}

func (r *PostgresRepo) FetchEmailByUUID(userGUID uuid.UUID) (string, error) {
	const method = "PostgresRepo.FetchEmailByUUID"

	row := r.pool.QueryRow(r.ctx, `SELECT email FROM users WHERE id=$1`, userGUID)
	var email string
	err := row.Scan(&email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", method, err)
	}
	return email, err
}

func (r *PostgresRepo) FetchUserByEmail(email string) (*users.User, error) {
	const method = "PostgresRepo.FetchUserByEmail"

	row := r.pool.QueryRow(r.ctx, `SELECT id, refresh_token FROM users WHERE email=$1`, email)
	var (
		userGUID         uuid.UUID
		refreshTokenHash string
	)
	err := row.Scan(&userGUID, &refreshTokenHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, users.ErrNoRefreshToken
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}

	return users.NewUser(userGUID, email, refreshTokenHash), nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}
