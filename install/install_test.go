package install

import (
	"errors"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	mocks "github.com/crowley-io/pack/testing"
	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c, "../testing/app.bin")
	o := dckOpts(t, c)

	wireMock(d, o, 0, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func TestInstallOnError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c, "../testing/app.bin")
	o := dckOpts(t, c)

	wireMock(d, o, 255, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestInstallWithDockerError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	e := errors.New("an error")
	setCnf(c, "../testing/app.bin")
	o := dckOpts(t, c)

	wireMock(d, o, 0, e)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, e, err)

}

func TestInstallWithConfigurationError(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c, "")

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, configuration.ErrOutputRequired, err)

}

func TestInstallWithNilConfiguration(t *testing.T) {

	d := &mocks.DockerMock{}

	err := Install(d, nil)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, configuration.ErrConfigurationEmpty, err)

}

func TestInstallWithNoOutput(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c, "file.txt")
	o := dckOpts(t, c)

	wireMock(d, o, 0, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestInstallDisabled(t *testing.T) {

	d := &mocks.DockerMock{}
	c := &configuration.Configuration{}
	setCnf(c, "../testing/app.bin")
	c.Install = configuration.Install{Disable: true}

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func wireMock(m *mocks.DockerMock, o docker.RunOptions, status int, err error) {
	m.On("Run", o, docker.NewLogStream()).Return(status, err)
}

func dckOpts(t *testing.T, c *configuration.Configuration) docker.RunOptions {

	e, err := GetEnv(c)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	v, err := GetVolumes(c)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	return docker.RunOptions{
		Image:   "debian",
		Command: "make",
		Env:     e,
		Volumes: v,
	}
}

func setCnf(c *configuration.Configuration, output string) {
	c.Install = configuration.Install{
		Command: "make",
		Path:    "/root",
		Image:   "debian",
		Output:  output,
	}
	c.Compose = configuration.Compose{
		Name: "debian",
	}
	c.Publish = configuration.Publish{
		Hostname: "localhost:5000",
	}
}
