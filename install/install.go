package install

import (
	"fmt"
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Install run compile instructions with a Docker container.
func Install(client docker.Docker, configuration *configuration.Configuration) error {

	if err := configuration.Validate(); err != nil {
		return err
	}

	if configuration.Install.Disable {
		return nil
	}

	env, err := GetEnv(configuration)

	if err != nil {
		return err
	}

	volumes, err := GetVolumes(configuration)

	if err != nil {
		return err
	}

	output := configuration.Install.Output
	image := configuration.Install.Image
	command := configuration.Install.Command

	option := docker.RunOptions{
		Image:   image,
		Command: command,
		Env:     env,
		Volumes: volumes,
	}

	exit, err := client.Run(option, docker.NewLogStream())

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

func pathExist(p string) bool {
	if _, err := os.Stat(p); err != nil {
		return false
	}
	return true
}
