package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWithPackerFile(t *testing.T) {

	c, err := Parse("../testing/packer.yml")

	if !assert.Nil(t, err) {
		t.Fatal("cannot parse file", err)
	}

	if assert.NotNil(t, c) {
		testConfArgs(t, c)
	}
}

func TestParseWithJsonFile(t *testing.T) {

	c, err := Parse("../testing/packer.json")

	if !assert.Nil(t, err) {
		t.Fatal("cannot parse file", err)
	}

	if assert.NotNil(t, c) {
		testConfArgs(t, c)
	}
}

func TestParseWithNoFile(t *testing.T) {

	c, err := Parse("../testing/nofile.yml")

	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func TestParseWithSyntaxError(t *testing.T) {

	c, err := Parse("../testing/packer-seq.yml")

	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func testConfArgs(t *testing.T, c *Configuration) {
	assert.Equal(t, "app.tar.gz", c.Output)
	assert.Equal(t, "packer-go", c.Install.Image)
	assert.Equal(t, "/usr/local/go", c.Install.Path)
	assert.Equal(t, "shen", c.Pack.Name)
}
