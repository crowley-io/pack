package install

import (
	"errors"
	"fmt"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

var (
	// ErrOutputRequired is returned when output isn't defined in the configuration.
	ErrOutputRequired = errors.New("configuration: output is required")
	// ErrPathRequired is returned when path isn't defined in the configuration.
	ErrPathRequired = errors.New("configuration: path is required")
	// ErrImageRequired is returned when path isn't defined in the configuration.
	ErrImageRequired = errors.New("configuration: image is required")
	// ErrConfigurationEmpty is returned when Install is called with an empty configuration.
	ErrConfigurationEmpty = errors.New("configuration is required")
)

// Install run compile instructions with a Docker container.
func Install(client docker.Docker, configuration *configuration.Configuration) error {

	if err := validateConfiguration(configuration); err != nil {
		return err
	}

	env, err := GetEnv(configuration)

	if err != nil {
		return err
	}

	volumes, err := GetVolumes(configuration)

	if err != nil {
		return err
	}

	image := configuration.Install.Image
	command := configuration.Install.Command

	option := docker.RunOptions{
		Image:   image,
		Command: command,
		Env:     env,
		Volumes: volumes,
	}

	exit, err := client.Run(option)

	if err != nil {
		return err
	}

	if exit != 0 {
		return fmt.Errorf("cannot run install: exit status %d", exit)
	}

	return nil
}

func validateConfiguration(c *configuration.Configuration) error {

	if c == nil {
		return ErrConfigurationEmpty
	}

	if c.Output == "" {
		return ErrOutputRequired
	}

	if c.Install.Path == "" {
		return ErrPathRequired
	}

	if c.Install.Image == "" {
		return ErrImageRequired
	}

	return nil
}
