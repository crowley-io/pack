package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultDockerEndpoint define the default docker endpoint used in Configuration
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

// Parse a file path and inflate a new Configuration
func Parse(path string) (*Configuration, error) {

	c := &Configuration{}
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)

	if err != nil {
		return nil, err
	}

	validate(c)

	return c, nil
}

func validate(c *Configuration) bool {

	if c.DockerEndpoint == "" {
		c.DockerEndpoint = DefaultDockerEndpoint
	}

	return true
}
