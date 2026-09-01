package api

import (
	"log/slog"
	"time"

	grpc "google.golang.org/grpc"
)

type ExecutionServer struct {
	UnimplementedExecutionServiceServer
	logger *slog.Logger
}

func (s *ExecutionServer) Execute(req *ExecuteRequest, stream grpc.ServerStreamingServer[ExecuteResponse]) error {
	for range 3 {
		stream.Send(&ExecuteResponse{Output: "streaming chunk\n"})
		time.Sleep(1 * time.Second)
	}
	return nil
}
func NewServer(log *slog.Logger) *ExecutionServer {
	return &ExecutionServer{
		logger: log,
	}
}
