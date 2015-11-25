package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

// Push options
const (
	pushRawJSONStream = false
)

// PushOptions contains the push configuration for the docker daemon.
type PushOptions struct {
	Name       string
	Repository string
	Registry   string
	Tag        string
}

// See Docker interface
func (d docker) Push(option PushOptions, stream LogStream) error {
	return d.client.PushImage(pushImageOptions(option, stream), pushAuthConfiguration(option))
}

func pushImageOptions(option PushOptions, stream LogStream) api.PushImageOptions {
	return api.PushImageOptions{
		Name:          option.Repository,
		Tag:           option.Tag,
		OutputStream:  stream.Out,
		RawJSONStream: pushRawJSONStream,
	}
}

func pushAuthConfiguration(option PushOptions) api.AuthConfiguration {
	return getAuthWithRegistry(option.Registry)
}
