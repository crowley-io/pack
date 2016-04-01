package configuration

import (
	"errors"
)

var (
	// ErrOutputRequired is returned when output isn't defined in the configuration.
	ErrOutputRequired = errors.New("configuration(output): a path is required")
	// ErrPathRequired is returned when path isn't defined in the configuration.
	ErrPathRequired = errors.New("configuration(install): path is required")
	// ErrImageRequired is returned when path isn't defined in the configuration.
	ErrImageRequired = errors.New("configuration(install): image is required")
	// ErrNameRequired is returned when image name isn't defined in the configuration.
	ErrNameRequired = errors.New("configuration(compose): name is required")
	// ErrHostnameRequired is returned when registry hostname isn't defined in the configuration.
	ErrHostnameRequired = errors.New("configuration(publish): hostname is required")
	// ErrConfigurationEmpty is returned when the configuration is empty.
	ErrConfigurationEmpty = errors.New("configuration is required")
)

// Validate return an error if the given Configuration has flaw.
func Validate(c *Configuration) error {

	if c == nil {
		return ErrConfigurationEmpty
	}

	if !c.Install.Disable {

		if c.Install.Output == "" {
			return ErrOutputRequired
		}

		if c.Install.Path == "" {
			return ErrPathRequired
		}

		if c.Install.Image == "" {
			return ErrImageRequired
		}

	}

	if c.Compose.Name == "" {
		return ErrNameRequired
	}

	if c.Publish.Hostname == "" {
		return ErrHostnameRequired
	}

	return nil
}
