package dockertest

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"os"
	"runtime"
)

func PrepareRunOptions(opts *dockertest.RunOptions, servicePorts ...docker.Port) *dockertest.RunOptions {
	if !IsDarwinOS() {
		return opts
	}

	if opts.PortBindings == nil {
		opts.PortBindings = make(map[docker.Port][]docker.PortBinding, len(servicePorts))
	}

	for _, port := range servicePorts {
		if _, ok := opts.PortBindings[port]; !ok {
			opts.PortBindings[port] = append(opts.PortBindings[port], docker.PortBinding{
				HostIP:   "0.0.0.0",
				HostPort: "0",
			})
		}
	}

	return opts
}

func IsDarwinOS() bool {
	return runtime.GOOS == "darwin"
}

func IsRunningInDockerContainer() bool {
	_, err := os.Stat("./.dockerenv")
	return err == nil
}

func Host(resource *dockertest.Resource) string {
	if IsDarwinOS() && !IsRunningInDockerContainer() {
		return "127.0.0.1"
	}

	return "0.0.0.0"
}
