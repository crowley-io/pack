package publish

import (
	"fmt"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Publish push the docker image into the docker registy.
func Publish(client docker.Docker, configuration *configuration.Configuration) error {

	if err := configuration.Validate(); err != nil {
		return err
	}

	name := configuration.Compose.Name
	hostname := configuration.Publish.Hostname
	stream := docker.NewLogStream()

	repository, tag := parseRepositoryTag(fmt.Sprintf("%s/%s", hostname, name))

	option := docker.PushOptions{
		Name:       name,
		Repository: repository,
		Registry:   hostname,
		Tag:        tag,
	}

	return client.Push(option, stream)

}
