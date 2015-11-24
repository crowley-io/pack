package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmptyConfiguration(t *testing.T) {

	c := New()

	if !assert.NotNil(t, c) {
		t.FailNow()
	}

	err := c.Validate()

	assert.NotNil(t, err)

}

func TestValidateNilConfiguration(t *testing.T) {

	err := Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, ErrConfigurationEmpty, err)

}

func TestValidateConfiguration(t *testing.T) {

	c := newTestConf()

	err := Validate(c)

	assert.Nil(t, err)

}

func TestValidateConfigurationEmptyOutput(t *testing.T) {

	c := newTestConf()
	c.Install.Output = ""

	err := Validate(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrOutputRequired, err)

}

func TestValidateConfigurationEmptyImage(t *testing.T) {

	c := newTestConf()
	c.Install.Image = ""

	err := Validate(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrImageRequired, err)

}

func TestValidateConfigurationEmptyPath(t *testing.T) {

	c := newTestConf()
	c.Install.Path = ""

	err := Validate(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrPathRequired, err)

}

func TestValidateConfigurationEmptyName(t *testing.T) {

	c := newTestConf()
	c.Compose.Name = ""

	err := Validate(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrNameRequired, err)

}

func TestValidateConfigurationHostname(t *testing.T) {

	c := newTestConf()
	c.Publish.Hostname = ""

	err := Validate(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrHostnameRequired, err)

}

func newTestConf() *Configuration {

	c := &Configuration{
		Install: Install{
			Image:  "debian",
			Path:   "/root",
			Output: "file",
		},
		Compose: Compose{
			Name: "app",
		},
		Publish: Publish{
			Hostname: "localhost:5000",
		},
	}

	values(c)
	return c
}
