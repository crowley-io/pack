package docker

import (
	"os"

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
func (d docker) Build(option BuildOptions) error {
	return d.client.BuildImage(buildImageOptions(option))
}

func buildImageOptions(option BuildOptions) api.BuildImageOptions {
	return api.BuildImageOptions{
		Name:                option.Name,
		NoCache:             option.NoCache,
		SuppressOutput:      quiet,
		RmTmpContainer:      rm,
		ForceRmTmpContainer: forceRm,
		Pull:                option.Pull,
		OutputStream:        os.Stdout,
		ContextDir:          option.Directory,
		AuthConfigs: api.AuthConfigurations{
			// TODO
			Configs: map[string]api.AuthConfiguration{
			// "quay.io": {
			// 	Username:      "foo",
			// 	Password:      "bar",
			// 	Email:         "baz",
			// 	ServerAddress: "quay.io",
			// },
			},
		},
	}
}
