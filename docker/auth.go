package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

func getAuth() *api.AuthConfigurations {

	c, err := api.NewAuthConfigurationsFromDockerCfg()

	if err != nil {
		return &api.AuthConfigurations{}
	}

	return c
}
