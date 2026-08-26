package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnandRaj2224/agentbox/internal/config"
	"github.com/AnandRaj2224/agentbox/internal/logger"
	"github.com/AnandRaj2224/agentbox/internal/orchestrator"
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

	cli, err := orchestrator.New(log)
	if err != nil {
		log.Error("error creating docker client", "Error", err)
		return
	}

	err = cli.PullImage(ctx, "alpine:latest")
	if err != nil {
		log.Error("error pulling the image", "Error", err)
		return

	}
	log.Info("image pulled", "image", "alpine:latest")

	ID, err := cli.CreateContainer(ctx, "alpine:latest", []string{"sleep", "10"})
	if err != nil {
		log.Error("error creating container", "Error", err)
		return

	}
	log.Info("container created", "id", ID)
	defer cli.RemoveContainer(context.Background(), ID)

	err = cli.StartContainer(ctx, ID)
	if err != nil {
		log.Error("error starting container", "Error", err)
		return

	}
	log.Info("container started", "id", ID)

	res,err := cli.InspectContainer(ctx,ID)
	if err != nil {
		log.Error("error inspecting container","Error",err)
		return
	}
	log.Info("inspection","state",res)

	err = cli.StopContainer(ctx, ID)
	if err != nil {
		return
	}
	
	err = cli.WaitContainer(ctx, ID)
	if err != nil {
		log.Error("error waiting for container", "Error", err)
		return

	}

	<-ctx.Done()
	log.Info("Server shutting down")
}
