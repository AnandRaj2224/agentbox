package orchestrator

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// hijackedIO wraps Docker's HijackedResponse into a standard io.ReadCloser.
type hijackedIO struct {
	io.Reader
	io.Closer
}

// DockerOrchestrator is a wrapper for Docker SDK client and logger.
type DockerOrchestrator struct {
	cli    *client.Client
	logger *slog.Logger
}

// CopyToContainer copies the contents of contentReader into a container
// at the specified destination path.
func (d *DockerOrchestrator) CopyToContainer(ctx context.Context, containerID string, destinationPath string, contentReader io.Reader) error {
	_, err := d.cli.CopyToContainer(ctx, containerID, client.CopyToContainerOptions{
		DestinationPath: destinationPath,
		Content:         contentReader,
	})
	if err != nil {
		return err
	}
	return nil
}

// AttachContainer opens a raw TCP socket between server and the Docker Daemon
// it is pulling the logs out of Docker and bringing them into the backend.
func (d *DockerOrchestrator) AttachContainer(ctx context.Context, containerID string) (io.ReadCloser, error) {
	resp, err := d.cli.ContainerAttach(ctx, containerID, client.ContainerAttachOptions{
		Stream: true, Stdout: true, Stderr: true, Logs: true,
	})
	if err != nil {
		return nil, err
	}

	return hijackedIO{
		Reader: resp.Reader,
		Closer: resp.Conn,
	}, nil
}

// PullImage downloads the container image from docker hub
func (d *DockerOrchestrator) PullImage(ctx context.Context, imageName string) error {
	reader, err := d.cli.ImagePull(
		ctx,
		imageName,
		client.ImagePullOptions{},
	)
	if err != nil {
		return err
	}

	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}
	return nil
}

// CreateContainer starts the creation of the docker container based on imageName and it returns a Unique ID to keep
// track of it.
func (d *DockerOrchestrator) CreateContainer(ctx context.Context, imageName string, cmd []string, cfg *container.HostConfig) (string, error) {
	resp, err := d.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: imageName,
		Config: &container.Config{
			Cmd: cmd,
		},
		HostConfig: cfg,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// StartContainer boots up the container process on the host OS.
func (d *DockerOrchestrator) StartContainer(ctx context.Context, containerID string) error {
	_, err := d.cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// WaitContainer blocks our Go application from moving forward until the container completely finishes running.
func (d *DockerOrchestrator) WaitContainer(ctx context.Context, containerID string) error {
	wait := d.cli.ContainerWait(ctx, containerID, client.ContainerWaitOptions{})
	select {
	case err := <-wait.Error:
		if err != nil {
			panic(err)
		}
	case <-wait.Result:
	}
	return nil
}

// RemoveContainer forcefully deletes the container and its volumes from the host machine.
func (d *DockerOrchestrator) RemoveContainer(ctx context.Context, containerID string) error {
	_, err := d.cli.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	if err != nil {
		return err
	}
	return nil
}

// InspectContainer queries the Docker API for the container's metadata.
func (d *DockerOrchestrator) InspectContainer(ctx context.Context, containerID string) (*container.State, error) {
	result, err := d.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, err
	}
	return result.Container.State, nil
}

// StopContainer sends a SIGTERM signal to the container, asking it to shut down gracefully.
func (d *DockerOrchestrator) StopContainer(ctx context.Context, containerID string) error {
	_, err := d.cli.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	if err != nil {
		return err
	}
	return nil
}

// New initializes and returns a new Docker Engine SDK client.
func New(log *slog.Logger) (*DockerOrchestrator, error) {
	cli, err := client.New()
	if err != nil {
		return nil, err
	}
	return &DockerOrchestrator{cli: cli, logger: log}, nil
}
