package compose

import (
	"errors"
	"os"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	mocks "github.com/crowley-io/pack/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c)

	o := dckOpts(t, c)
	if !assert.NotNil(t, o) {
		t.FailNow()
	}

	ls := docker.NewLogStream()

	wireMock(d, ls, o, nil)

	err := Compose(d, ls, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func TestComposeWithConfigurationError(t *testing.T) {

	c := &configuration.Configuration{}
	d := &mocks.DockerMock{}
	ls := docker.NewLogStream()

	err := Compose(d, ls, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestComposeOnError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	e := errors.New("an error")
	setCnf(c)

	o := dckOpts(t, c)
	if !assert.NotNil(t, o) {
		t.FailNow()
	}

	ls := docker.NewLogStream()

	wireMock(d, ls, o, e)

	err := Compose(d, ls, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, e, err)

}

func wireMock(m *mocks.DockerMock, ls docker.LogStream, o docker.BuildOptions, err error) {
	m.On("Build", o, ls).Return(err)
	m.On("ImageID", o.Name).Return("26913aba19ca").Twice()
}

func dckOpts(t *testing.T, c *configuration.Configuration) docker.BuildOptions {

	pwd, err := os.Getwd()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	return docker.BuildOptions{
		Name:      c.Compose.Name,
		Directory: pwd,
		Pull:      false,
		NoCache:   false,
	}
}

func setCnf(c *configuration.Configuration) {
	c.Install = configuration.Install{
		Command: "make",
		Path:    "/root",
		Image:   "debian",
		Output:  "libapp.so",
	}
	c.Compose = configuration.Compose{
		Name: "debian",
	}
	c.Publish = configuration.Publish{
		Hostname: "localhost:5000",
	}
}
