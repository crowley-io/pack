package publish

import (
	"errors"
	"testing"

	parser "github.com/crowley-io/docker-parser"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	mocks "github.com/crowley-io/pack/testing"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c)
	ls := docker.NewLogStream()
	topts, popts := options(c)

	wireMock(d, ls, topts, nil, popts, nil)

	err := Publish(d, ls, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func TestPublishWithConfigurationError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	ls := docker.NewLogStream()

	err := Publish(d, ls, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestPublishOnTagError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	e := errors.New("an error")
	setCnf(c)
	ls := docker.NewLogStream()
	topts, popts := options(c)

	wireMock(d, ls, topts, e, popts, nil)

	err := Publish(d, ls, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, e, err)

}

func TestName(t *testing.T) {
	assertConfParsing(t, "debian:8.2", func(p *parser.Reference) string {
		return p.Name()
	})
}

func TestRegistry(t *testing.T) {
	assertConfParsing(t, "localhost:5000", func(p *parser.Reference) string {
		return p.Registry()
	})
}

func TestRemote(t *testing.T) {
	assertConfParsing(t, "localhost:5000/debian:8.2", func(p *parser.Reference) string {
		return p.Remote()
	})
}

func TestRepository(t *testing.T) {
	assertConfParsing(t, "localhost:5000/debian", func(p *parser.Reference) string {
		return p.Repository()
	})
}

func TestTag(t *testing.T) {
	assertConfParsing(t, "8.2", func(p *parser.Reference) string {
		return p.Tag()
	})
}

func assertConfParsing(t *testing.T, expected string, callback func(*parser.Reference) string) {

	c := &configuration.Configuration{}
	setCnf(c)

	s := remote(c)
	p, err := parser.Parse(s)

	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expected, callback(p))
}

func wireMock(m *mocks.DockerMock, ls docker.LogStream, topts docker.TagOptions, terr error, popts docker.PushOptions, perr error) {
	m.On("Tag", topts).Return(terr)
	if terr == nil {
		m.On("Push", popts, ls).Return(perr)
		m.On("RemoveImage", "localhost:5000/debian:8.2").Return(nil)
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

func options(configuration *configuration.Configuration) (tagOpts docker.TagOptions, pushOpts docker.PushOptions) {

	if reference, err := parser.Parse(remote(configuration)); err == nil {

		tagOpts = docker.TagOptions{
			Name:       reference.Name(),
			Repository: reference.Repository(),
			Tag:        reference.Tag(),
		}

		pushOpts = docker.PushOptions{
			Name:       reference.Name(),
			Repository: reference.Repository(),
			Registry:   reference.Registry(),
			Tag:        reference.Tag(),
		}
	}

	return
}
