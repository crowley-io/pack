package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

// Tag options
const (
	tagForce = true
)

// TagOptions contains the tag configuration for the docker daemon.
type TagOptions struct {
	Name       string
	Repository string
	Tag        string
}

// See Docker interface
func (d docker) Tag(option TagOptions) error {
	return d.client.TagImage(option.Name, tagImageOptions(option))
}

func tagImageOptions(option TagOptions) api.TagImageOptions {
	return api.TagImageOptions{
		Repo:  option.Repository,
		Tag:   option.Tag,
		Force: tagForce,
	}
}
