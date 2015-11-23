package configuration

const (
	// DefaultDockerEndpoint define the default docker endpoint used in Configuration
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

// Values define and use the default Values for the given
// Configuration, which will be override later by user's configuration.
func values(c *Configuration) {

	c.DockerEndpoint = DefaultDockerEndpoint
	c.Compose.Pull = true

}
