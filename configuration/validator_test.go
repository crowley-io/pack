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
	c.Output = ""

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

func newTestConf() *Configuration {

	c := &Configuration{
		Output: "file",
		Install: Install{
			Image: "debian",
			Path:  "/root",
		},
		Compose: Compose{
			Name: "app",
		},
	}

	values(c)
	return c
}
