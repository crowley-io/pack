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
	Name string
}
