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
	registry := configuration.Publish.Hostname
	repository, tag := parseRepositoryTag(fmt.Sprintf("%s/%s", registry, name))

	to := docker.TagOptions{
		Name:       name,
		Repository: repository,
		Tag:        tag,
	}

	po := docker.PushOptions{
		Name:       name,
		Repository: repository,
		Registry:   registry,
		Tag:        tag,
	}

	if err := client.RemoveImage(repository); err != nil {
		return err
	}

	if err := client.Tag(to); err != nil {
		return err
	}

	return client.Push(po, docker.NewLogStream())

}
