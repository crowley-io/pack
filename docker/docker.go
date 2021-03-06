package docker

import (
	"github.com/crowley-io/pack/configuration"
	api "github.com/fsouza/go-dockerclient"
)

// Docker wrap underlaying docker client to expose only required functions.
type Docker interface {
	// Ping pings the docker server
	Ping() error
	// Logs attach a stream on a running container to read stdout and stderr
	// output from docker logs.
	Logs(id string, stream LogStream) error
	// Run create and start a container to execute a runnable.
	// Return the exit code of the container status, an error otherwise.
	Run(option RunOptions, stream LogStream) (int, error)
	// Build create a new image from a Dockerfile.
	Build(option BuildOptions, stream LogStream) error
	// Tag adds a tag to the image for a repository.
	Tag(option TagOptions) error
	// Push pushes an image to a remote registry.
	Push(option PushOptions, stream LogStream) error
	// ImageID returns an image ID by its name.
	ImageID(name string) string
	// RemoveImage removes an image by its name or ID.
	RemoveImage(name string) error
}

// The default implementation of the Docker interface.
type docker struct {
	client *api.Client
}

// See Docker interface
func (d docker) Ping() error {
	return d.client.Ping()
}

// New return a Docker client
func New(configuration *configuration.Configuration) (Docker, error) {

	c, err := api.NewClient(configuration.DockerEndpoint)

	if err != nil {
		return nil, err
	}

	d := &docker{}
	d.client = c

	if err = d.Ping(); err != nil {
		return nil, err
	}

	return d, nil

}
