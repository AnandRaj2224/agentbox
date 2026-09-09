package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ExecutionRecord struct {
	ID        uuid.UUID
	Runtime   string
	Timestamp time.Time
	ExitCode  int
	Duration  time.Duration
	Code      string
	Output    string
}

type Repository interface {
	SaveExecution(ctx context.Context, record ExecutionRecord) error
}
