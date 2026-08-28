package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnandRaj2224/agentbox/internal/api"
	"github.com/AnandRaj2224/agentbox/internal/config"
	"github.com/AnandRaj2224/agentbox/internal/logger"
	"google.golang.org/grpc"
)

func main() {
	config := config.Load()
	ctx := context.Background()
	log := logger.New()

	log.Info("server started",
		"port", config.Port,
		"environment", config.Env,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
	)
	defer stop()

	listener, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Error("error listen on port", config.Port, err)
		return
	}

	grpcServer := grpc.NewServer()
	api.RegisterExecutionServiceServer(grpcServer, api.NewServer(log))
	go grpcServer.Serve(listener)

	<-ctx.Done()
	grpcServer.GracefulStop()
	log.Info("Server shutting down")
}
