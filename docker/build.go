package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

// Build options
const (
	buildQuiet   = false
	buildRm      = true
	buildForceRm = true
	buildMemory  = 0
	buildMemswap = 0
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
		SuppressOutput:      buildQuiet,
		RmTmpContainer:      buildRm,
		ForceRmTmpContainer: buildForceRm,
		Pull:                option.Pull,
		OutputStream:        stream.Decoder,
		ContextDir:          option.Directory,
		RawJSONStream:       rawJSONStream,
		Memory:              buildMemory,
		Memswap:             buildMemswap,
	}

	if auth := getAuth(); auth != nil {
		opts.AuthConfigs = *auth
	}

	return opts
}
