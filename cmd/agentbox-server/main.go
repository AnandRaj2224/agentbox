package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnandRaj2224/agentbox/internal/config"
	"github.com/AnandRaj2224/agentbox/internal/logger"
)

func main() {
	config := config.Load()

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

	<-ctx.Done()

	log.Info("Server shutting down")
}
