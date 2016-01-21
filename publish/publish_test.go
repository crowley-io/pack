package publish

import (
	"errors"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	mocks "github.com/crowley-io/pack/testing"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c)
	topts, popts := options(c)

	wireMock(d, topts, nil, popts, nil)

	err := Publish(d, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func TestPublishWithConfigurationError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}

	err := Publish(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestPublishOnTagError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	e := errors.New("an error")
	setCnf(c)
	topts, popts := options(c)

	wireMock(d, topts, e, popts, nil)

	err := Publish(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, e, err)

}

func TestName(t *testing.T) {

	e := "debian:8.2"
	c := &configuration.Configuration{}
	setCnf(c)

	i := Name(c)

	assert.Equal(t, e, i)

}

func TestRegistry(t *testing.T) {

	e := "localhost:5000"
	c := &configuration.Configuration{}
	setCnf(c)

	i := Registry(c)

	assert.Equal(t, e, i)

}

func TestRemote(t *testing.T) {

	e := "localhost:5000/debian:8.2"
	c := &configuration.Configuration{}
	setCnf(c)

	i := Remote(c)

	assert.Equal(t, e, i)

}

func TestRepository(t *testing.T) {

	e := "localhost:5000/debian"
	c := &configuration.Configuration{}
	setCnf(c)

	i := Repository(c)

	assert.Equal(t, e, i)

}

func TestTag(t *testing.T) {

	e := "8.2"
	c := &configuration.Configuration{}
	setCnf(c)

	i := Tag(c)

	assert.Equal(t, e, i)

}

func wireMock(m *mocks.DockerMock, topts docker.TagOptions, terr error, popts docker.PushOptions, perr error) {
	m.On("Tag", topts).Return(terr)
	if terr == nil {
		m.On("Push", popts, docker.NewLogStream()).Return(perr)
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
		Name: "debian:8.2",
	}
	c.Publish = configuration.Publish{
		Hostname: "localhost:5000",
	}
}
