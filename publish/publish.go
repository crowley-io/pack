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

	reference, err := parser.Parse(remote(configuration))

	if err != nil {
		return err
	}

	tagOpts := docker.TagOptions{
		Name:       reference.Name(),
		Repository: reference.Repository(),
		Tag:        reference.Tag(),
	}

	pushOpts := docker.PushOptions{
		Name:       reference.Name(),
		Repository: reference.Repository(),
		Registry:   reference.Registry(),
		Tag:        reference.Tag(),
	}

	if err = client.Tag(tagOpts); err != nil {
		return err
	}

	stream := docker.NewLogStream()
	err = client.Push(pushOpts, stream)

	// Remove registry tag
	if err2 := client.RemoveImage(reference.Remote()); err2 != nil && err == nil {
		err = err2
	}

	if err3 := stream.Close(); err3 != nil && err == nil {
		err = err3
	}

	return err
}

// Remote returns the image's remote identifier. (ie: registry/name[:tag])
func remote(configuration *configuration.Configuration) string {
	return fmt.Sprintf("%s/%s", configuration.Publish.Hostname, configuration.Compose.Name)
}
