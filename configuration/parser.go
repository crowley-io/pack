package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Parse a file path and inflate a new Configuration
func Parse(path string) (*Configuration, error) {

	c := &Configuration{}
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}

	if err = Validate(c); err != nil {
		return nil, err
	}

	return c, nil
}
