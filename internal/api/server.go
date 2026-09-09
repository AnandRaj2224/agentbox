package api

import (
	context "context"
	"log/slog"

	"github.com/AnandRaj2224/agentbox/internal/orchestrator"
	"github.com/AnandRaj2224/agentbox/internal/runtime"
	"github.com/AnandRaj2224/agentbox/internal/sandbox"
	"github.com/moby/moby/api/pkg/stdcopy"
	grpc "google.golang.org/grpc"
)

// grpcWriter intercept raw bytes and wraps them in a Protobuf struct
// then push it down the network stream.
type grpcWriter struct {
	stream grpc.ServerStreamingServer[ExecuteResponse]
}

// Write intercept bytes wrap them in grpc ExecuteResponse struct, and calls w.stream.Send().
func (w *grpcWriter) Write(p []byte) (int, error) {
	w.stream.Send(&ExecuteResponse{Output: string(p)})
	return len(p), nil
}

// ExecutionServer is a container for dependencies
// giving every incoming request has access to the tools it needs.
type ExecutionServer struct {
	UnimplementedExecutionServiceServer
	logger *slog.Logger
	cli    *orchestrator.DockerOrchestrator
}

// Execute is a Request Handler that reads users's request.
// ask for specific runtime,tells orchestrator package to spin up the container and attach to the logs.
// It tells stdcopy to pump the logs into the network stream.
// When the container dies, this method finishes, the network stream closes, and the container is deleted.
func (s *ExecutionServer) Execute(req *ExecuteRequest, stream grpc.ServerStreamingServer[ExecuteResponse]) error {
	ctx := stream.Context()
	cfg := sandbox.SandboxConfig{MemoryMB: 50, CPULimit: 0.5}
	rt, err := runtime.GetRuntime(req.Runtime)
	if err != nil {
		return err
	}

	err = s.cli.PullImage(ctx, rt.Image())
	if err != nil {
		return err
	}
	containerID, err := s.cli.CreateContainer(ctx, rt.Image(), rt.Command("main"+rt.Extension()), sandbox.NewHostConfig(cfg))
	if err != nil {
		return err
	}
	defer s.cli.RemoveContainer(context.Background(), containerID)

	tarReader, err := sandbox.CreateTarArchive("main"+rt.Extension(), req.SourceCode)
	err = s.cli.CopyToContainer(ctx, containerID, "/", tarReader)

	logReader, err := s.cli.AttachContainer(ctx, containerID)
	if err != nil {
		return err
	}
	defer logReader.Close()

	err = s.cli.StartContainer(ctx, containerID)
	if err != nil {
		return err
	}

	writer := &grpcWriter{stream: stream}
	stdcopy.StdCopy(writer, writer, logReader)

	err = s.cli.WaitContainer(ctx, containerID)
	if err != nil {
		return err
	}
	return nil
}

// NewServer is a Constructor for *ExecutionServer.
func NewServer(log *slog.Logger, cli *orchestrator.DockerOrchestrator) *ExecutionServer {
	return &ExecutionServer{
		logger: log,
		cli:    cli,
	}
}
