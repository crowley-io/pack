package install

import (
	"errors"
	"fmt"
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

var (
	// ErrOutputRequired is returned when output isn't defined in the configuration.
	ErrOutputRequired = errors.New("configuration(install): output is required")
	// ErrPathRequired is returned when path isn't defined in the configuration.
	ErrPathRequired = errors.New("configuration(install): path is required")
	// ErrImageRequired is returned when path isn't defined in the configuration.
	ErrImageRequired = errors.New("configuration(install): image is required")
	// ErrConfigurationEmpty is returned when Install is called with an empty configuration.
	ErrConfigurationEmpty = errors.New("configuration is required")
)

// Install run compile instructions with a Docker container.
func Install(client docker.Docker, configuration *configuration.Configuration) error {

	if err := ValidateConfiguration(configuration); err != nil {
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

	output := configuration.Output
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

	if !pathExist(output) {
		return fmt.Errorf("file not found: %s", output)
	}

	return nil
}

// ValidateConfiguration return an error if the given Configuration has flaw.
func ValidateConfiguration(c *configuration.Configuration) error {

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

func pathExist(p string) bool {
	if _, err := os.Stat(p); err != nil {
		return false
	}
	return true
}
