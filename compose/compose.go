package compose

import (
	//	"errors"
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Compose create a new image.
func Compose(client docker.Docker, configuration *configuration.Configuration) error {

	if err := configuration.Validate(); err != nil {
		return err
	}

	name := configuration.Compose.Name
	noCache := configuration.Compose.NoCache
	pull := configuration.Compose.Pull
	directory, err := os.Getwd()

	if err != nil {
		return err
	}

	id := client.ImageID(name)
	defer client.RemoveImage(id)

	option := docker.BuildOptions{
		Name:      name,
		Directory: directory,
		Pull:      pull,
		NoCache:   noCache,
	}

	return client.Build(option, docker.NewLogStream())

}
