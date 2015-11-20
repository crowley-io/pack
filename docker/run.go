package docker

import (
	"fmt"
	"os"
	"strings"

	api "github.com/fsouza/go-dockerclient"
)

// Container options
const (
	containerName   = ""
	useTTy          = false
	attachStdout    = true
	attachStderr    = true
	hostNetworkMode = "bridge"
	networkDisabled = false
	removeVolumes   = true
	forceRemove     = false
)

// RunOptions contains the run configuration of the docker container
type RunOptions struct {
	Image   string
	Command string
	Env     []string
	Volumes []string
}

// Run create and start a container to execute a runnable.
// Return the exit code of the container status.
func (d Docker) Run(option RunOptions) (int, error) {

	e, err := d.client.CreateContainer(createContainerOptions(option))

	if err != nil {
		return 0, err
	}

	id := e.ID

	if err = d.client.StartContainer(id, nil); err != nil {
		return 0, err
	}

	err = d.Logs(id, LogStream{Out: os.Stdout, Err: os.Stderr})

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	exit, err := d.client.WaitContainer(id)

	if err != nil {
		return 0, err
	}

	if err = d.client.RemoveContainer(removeContainerOptions(id)); err != nil {
		return 0, err
	}

	return exit, nil

}

func createContainerOptions(option RunOptions) api.CreateContainerOptions {
	return api.CreateContainerOptions{
		Name: containerName,
		Config: &api.Config{
			AttachStdout:    attachStdout,
			AttachStderr:    attachStderr,
			Tty:             useTTy,
			Env:             option.Env,
			NetworkDisabled: networkDisabled,
			Image:           option.Image,
			Cmd:             strings.Fields(option.Command),
		},
		HostConfig: &api.HostConfig{
			NetworkMode: hostNetworkMode,
			Binds:       option.Volumes,
		},
	}
}

func removeContainerOptions(id string) api.RemoveContainerOptions {
	return api.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: removeVolumes,
		Force:         forceRemove,
	}
}
