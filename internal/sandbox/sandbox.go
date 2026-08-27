package sandbox

import "github.com/moby/moby/api/types/container"

// SandboxConfig creates a clean API.
// example -> 5 MB instead of 5242880 bytes(docker format).
type SandboxConfig struct {
	MemoryMB int
	CPULimit float64
}

// NewHostConfig acts as a Adapter that translates human readable
// config to exact specifications docker SDK requires.
func NewHostConfig(cfg SandboxConfig) *container.HostConfig {
	return &container.HostConfig{
		Resources: container.Resources{
			Memory:   int64(cfg.MemoryMB * 1024 * 1024),
			NanoCPUs: int64(cfg.CPULimit * 1e9),
		},
		// Disable networking to isolate the sandbox from external networks.
		NetworkMode: container.NetworkMode("none"),
	}
}
