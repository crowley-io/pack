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

	t, p := options(configuration)

	if err := client.Tag(t); err != nil {
		return err
	}

	return client.Push(p, docker.NewLogStream())
}

func options(configuration *configuration.Configuration) (docker.TagOptions, docker.PushOptions) {

	name := Name(configuration)
	registry := Registry(configuration)
	repository, tag := parseRepositoryTag(Remote(configuration))

	t := docker.TagOptions{
		Name:       name,
		Repository: repository,
		Tag:        tag,
	}

	p := docker.PushOptions{
		Name:       name,
		Repository: repository,
		Registry:   registry,
		Tag:        tag,
	}

	return t, p
}

// Name returns the image's name. (ie: name[:tag])
func Name(configuration *configuration.Configuration) string {
	return configuration.Compose.Name
}

// Registry returns the image's registry. (ie: host[:port])
func Registry(configuration *configuration.Configuration) string {
	return configuration.Publish.Hostname
}

// Remote returns the image's remote identifier. (ie: registry/name[:tag])
func Remote(configuration *configuration.Configuration) string {
	return fmt.Sprintf("%s/%s", Registry(configuration), Name(configuration))
}

// Repository returns the image's repository. (ie: registry/name)
func Repository(configuration *configuration.Configuration) string {
	r, _ := parseRepositoryTag(Remote(configuration))
	return r
}

// Tag returns the image's tag.
func Tag(configuration *configuration.Configuration) string {
	_, t := parseRepositoryTag(Remote(configuration))
	return t
}
