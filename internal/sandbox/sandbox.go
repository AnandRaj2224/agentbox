package sandbox

import "github.com/moby/moby/api/types/container"

type SandboxConfig struct {
	MemoryMB int
	CPULimit float64
}

func NewHostConfig(cfg SandboxConfig) *container.HostConfig {
	return &container.HostConfig{
		Resources: container.Resources{
			Memory:   int64(cfg.MemoryMB * 1024 * 1024),
			NanoCPUs: int64(cfg.CPULimit * 1e9),
		},
		NetworkMode: container.NetworkMode("none"),
	}
}
