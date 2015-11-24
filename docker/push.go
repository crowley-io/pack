package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

const (
	pushRawJSONStream = false
)

// PushOptions contains the push configuration for the docker daemon.
type PushOptions struct {
	Name     string
	Tag      string
	Registry string
	Command  string
	Env      []string
	Volumes  []string
}

// See Docker interface
func (d docker) Push(option PushOptions, stream LogStream) error {
	return d.client.PushImage(pushImageOptions(option, stream), pushAuthConfiguration(option))
}

func pushImageOptions(option PushOptions, stream LogStream) api.PushImageOptions {
	return api.PushImageOptions{
		Name:          option.Name,
		Tag:           option.Tag,
		Registry:      option.Registry,
		OutputStream:  stream.Out,
		RawJSONStream: pushRawJSONStream,
	}
}

func pushAuthConfiguration(option PushOptions) api.AuthConfiguration {
	// TODO
	return api.AuthConfiguration{}
}
