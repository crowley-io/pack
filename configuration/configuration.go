package configuration

// Configuration contains pack runtime instructions.
type Configuration struct {
	DockerEndpoint string
	Output         string
	Install        Install
	Pack           Pack
}

// Install contains install runtime's configuration.
type Install struct {
	Image       string
	Path        string
	Command     string
	Environment []string
	Volumes     []string
}

// Pack contains pack runtime's configuration.
type Pack struct {
	Name string
}
