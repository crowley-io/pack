package docker

import (
	"strings"

	api "github.com/fsouza/go-dockerclient"
)

func getAuth() *api.AuthConfigurations {

	c, err := api.NewAuthConfigurationsFromDockerCfg()

	if err != nil {
		return &api.AuthConfigurations{}
	}

	return c
}

func getAuthWithRegistry(registry string) api.AuthConfiguration {

	auth := getAuth()

	if auth != nil {
		for k, v := range auth.Configs {
			if indexName(k) == indexName(registry) {
				return v
			}
		}
	}

	return api.AuthConfiguration{}
}

func getAuthWithImage(image string) api.AuthConfiguration {

	empty := api.AuthConfiguration{}
	auth := getAuthWithRegistry(image)

	if auth == empty {
		return getAuthWithRegistry("docker.io")
	}

	return auth
}

func indexName(val string) string {

	val = toHostname(val)

	// 'index.docker.io' => 'docker.io'
	if val == "index.docker.io" {
		return "docker.io"
	}

	return val
}

func toHostname(url string) string {

	s := url

	if strings.HasPrefix(url, "http://") {
		s = strings.Replace(url, "http://", "", 1)
	} else if strings.HasPrefix(url, "https://") {
		s = strings.Replace(url, "https://", "", 1)
	}

	p := strings.SplitN(s, "/", 2)

	return p[0]
}
