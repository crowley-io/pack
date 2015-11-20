package docker

import (
	"github.com/crowley-io/pack/configuration"
	api "github.com/fsouza/go-dockerclient"
)

// Docker wrap underlaying docker client to expose only required functions.
type Docker struct {
	client *api.Client
}

// New return a Docker client
func New(configuration configuration.Configuration) (*Docker, error) {

	c, err := api.NewClient(configuration.DockerEndpoint)

	if err != nil {
		return nil, err
	}

	d := &Docker{}
	d.client = c

	if err = d.Ping(); err != nil {
		return nil, err
	}

	return d, nil
}

// Ping pings the docker server
func (d Docker) Ping() error {
	return d.client.Ping()
}
