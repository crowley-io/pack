package docker

import (
	"fmt"
	"strings"

	parser "github.com/crowley-io/docker-parser"
	api "github.com/fsouza/go-dockerclient"
)

// Run options
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
	Links   []string
}

// See Docker interface
func (d docker) Run(option RunOptions, stream LogStream) (int, error) {

	r, err := parser.Parse(option.Image)
	if err != nil {
		return 0, err
	}

	err = d.client.PullImage(pullImageOptions(r.Remote(), stream), pullAuthConfiguration(option))
	if err != nil {
		return 0, err
	}

	e, err := d.client.CreateContainer(createContainerOptions(option))

	if err != nil {
		return 0, err
	}

	id := e.ID

	if err = d.client.StartContainer(id, nil); err != nil {
		return 0, err
	}

	err = d.Logs(id, stream)

	if err != nil {
		fmt.Fprint(stream.Err, err)
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

func pullAuthConfiguration(option RunOptions) api.AuthConfiguration {
	return getAuthWithImage(option.Image)
}

func pullImageOptions(remote string, stream LogStream) api.PullImageOptions {
	return api.PullImageOptions{
		Repository:    remote,
		OutputStream:  stream.Decoder,
		RawJSONStream: rawJSONStream,
	}
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
			Links:       option.Links,
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
