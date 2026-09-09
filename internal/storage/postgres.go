package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func (s *PostgresStore) SaveExecution(ctx context.Context, record ExecutionRecord) error {
	query := "INSERT INTO executions (execution_id, runtime, timestamp, exitcode, execution_duration, source_code, program_output) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := s.pool.Exec(ctx, query, record.ID, record.Runtime, record.Timestamp, record.ExitCode, record.Duration.Milliseconds(), record.Code, record.Output)

	return err
}

func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}
