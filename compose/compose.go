package compose

import (
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Compose create a new image.
func Compose(client docker.Docker, stream docker.LogStream, configuration *configuration.Configuration) error {

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

	option := docker.BuildOptions{
		Name:      name,
		Directory: directory,
		Pull:      pull,
		NoCache:   noCache,
	}

	err = client.Build(option, stream)

	if newid := client.ImageID(name); newid != id {
		// Remove previous image since id doesn't match.
		if err2 := client.RemoveImage(id); err == nil {
			err = err2
		}
	}

	return err
}
