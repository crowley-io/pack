package configuration

// Configuration contains pack runtime instructions.
type Configuration struct {
	DockerEndpoint string
	Output         string
	Install        Install
	Compose        Compose
}

// Install contains install runtime's configuration.
type Install struct {
	Image       string
	Path        string
	Command     string
	Environment []string
	Volumes     []string
}

// Compose contains pack runtime's configuration.
type Compose struct {
	Name    string
	NoCache bool `yaml:"no-cache"`
	Pull    bool
}

// Validate is just a wrapper for the static function with the same name.
// Which mean that it will return an error if the given Configuration has flaw.
func (c *Configuration) Validate() error {
	return Validate(c)
}

// New return a default Configuration.
func New() *Configuration {
	c := &Configuration{}
	values(c)
	return c
}
