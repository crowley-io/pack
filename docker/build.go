package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

// Build options
const (
	quiet   = false
	rm      = true
	forceRm = true
)

// BuildOptions contains the build configuration for the docker daemon.
type BuildOptions struct {
	Name      string
	Directory string
	Pull      bool
	NoCache   bool
}

// See Docker interface
func (d docker) Build(option BuildOptions, stream LogStream) error {
	return d.client.BuildImage(buildImageOptions(option, stream))
}

func buildImageOptions(option BuildOptions, stream LogStream) api.BuildImageOptions {

	opts := api.BuildImageOptions{
		Name:                option.Name,
		NoCache:             option.NoCache,
		SuppressOutput:      quiet,
		RmTmpContainer:      rm,
		ForceRmTmpContainer: forceRm,
		Pull:                option.Pull,
		OutputStream:        stream.Out,
		ContextDir:          option.Directory,
	}

	if auth := getAuth(); auth != nil {
		opts.AuthConfigs = *auth
	}

	return opts
}
