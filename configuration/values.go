package configuration

const (
	// DefaultDockerEndpoint define the default docker endpoint used in Configuration
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

// Values define and use the default Values for the given
// Configuration, if required.
func Values(c *Configuration) {

	if c.DockerEndpoint == "" {
		c.DockerEndpoint = DefaultDockerEndpoint
	}

}
