package publish

import (
	"fmt"

	parser "github.com/crowley-io/docker-parser"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Publish push the docker image into the docker registy.
func Publish(client docker.Docker, configuration *configuration.Configuration) error {

	if err := configuration.Validate(); err != nil {
		return err
	}

	t, p, err := options(configuration)

	if err != nil {
		return err
	}

	if err = client.Tag(t); err != nil {
		return err
	}

	stream := docker.NewLogStream()
	err = client.Push(p, stream)

	// Remove registry tag
	if err2 := client.RemoveImage(remote(configuration)); err2 != nil && err == nil {
		err = err2
	}

	if err3 := stream.Close(); err3 != nil && err == nil {
		err = err3
	}

	return err
}

func parse(configuration *configuration.Configuration) (name, registry, repository, tag string, err error) {

	remote := remote(configuration)

	if name, err = parser.Name(remote); err != nil {
		return
	}

	if registry, err = parser.Registry(remote); err != nil {
		return
	}

	if repository, err = parser.Repository(remote); err != nil {
		return
	}

	if tag, err = parser.Tag(remote); err != nil {
		return
	}

	return
}

func options(configuration *configuration.Configuration) (docker.TagOptions, docker.PushOptions, error) {

	name, registry, repository, tag, err := parse(configuration)

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

	return t, p, err
}

// Remote returns the image's remote identifier. (ie: registry/name[:tag])
func remote(configuration *configuration.Configuration) string {
	return fmt.Sprintf("%s/%s", configuration.Publish.Hostname, configuration.Compose.Name)
}
