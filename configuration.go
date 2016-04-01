package main

import (
	"github.com/crowley-io/pack/configuration"
)

func createFileConfigurationParser(path *string) func() (*configuration.Configuration, error) {
	return func() (*configuration.Configuration, error) {
		return configuration.Parse(*path)
	}
}

func createRemoteConfigurationParser(remote *string,
	nocache, nopull *bool) func() (*configuration.Configuration, error) {

	return func() (*configuration.Configuration, error) {

		c := configuration.New()
		c.DisableInstall()

		if err := c.Configure(*remote); err != nil {
			return nil, err
		}

		if *nocache {
			c.DisableCache()
		}

		if *nopull {
			c.DisablePull()
		}

		return c, nil
	}
}
