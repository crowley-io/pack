package configuration

import (
	parser "github.com/crowley-io/docker-parser"
)

// Configuration contains pack runtime instructions.
type Configuration struct {
	DockerEndpoint string `yaml:"docker-endpoint"`
	Install        Install
	Compose        Compose
	Publish        Publish
}

// Install contains install runtime's configuration.
type Install struct {
	Disable     bool
	Image       string
	Path        string
	Output      string
	Command     string
	Environment []string
	Volumes     []string
	Links       []string
}

// Compose contains compose runtime's configuration.
type Compose struct {
	Name    string
	NoCache bool `yaml:"no-cache"`
	Pull    bool
}

// Publish contains publish runtime's configuration.
type Publish struct {
	Hostname string
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

// Configure will analyze the given remote identifier and define the required parameters for
// Compose and Publish modules.
func (c *Configuration) Configure(remote string) error {

	r, err := parser.Parse(remote)
	if err != nil {
		return err
	}

	c.Compose.Name = r.Name()
	c.Publish.Hostname = r.Registry()
	return nil
}

// EnableInstall enable Install module.
func (c *Configuration) EnableInstall() {
	c.Install.Disable = false
}

// DisableInstall disable Install module.
func (c *Configuration) DisableInstall() {
	c.Install.Disable = true
}

// EnableCache enable cache for Compose module.
func (c *Configuration) EnableCache() {
	c.Compose.NoCache = false
}

// DisableCache disable cache for Compose module.
func (c *Configuration) DisableCache() {
	c.Compose.NoCache = true
}

// EnablePull enable pull for Compose module.
func (c *Configuration) EnablePull() {
	c.Compose.Pull = true
}

// DisablePull disable pull for Compose module.
func (c *Configuration) DisablePull() {
	c.Compose.Pull = false
}
