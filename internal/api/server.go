package api

import (
	"context"
	"log/slog"
)

type ExecutionServer struct {
	UnimplementedExecutionServiceServer
	logger *slog.Logger
}

func (s *ExecutionServer) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	s.logger.Info("dummy log", "sourceCode", req.SourceCode, "runtime", req.Runtime)
	return &ExecuteResponse{
		Output:   "received!",
		ExitCode: 0,
	}, nil
}

func NewServer(log *slog.Logger) *ExecutionServer {
	return &ExecutionServer{
		logger: log,
	}
}
