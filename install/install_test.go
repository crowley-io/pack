package install

import (
	"errors"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateConfigurationEmptyOutput(t *testing.T) {

	c := &configuration.Configuration{
		Install: configuration.Install{
			Path:  "/root",
			Image: "debian",
		},
	}
	err := validateConfiguration(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrOutputRequired, err)

}

func TestValidateConfigurationEmptyImage(t *testing.T) {

	c := &configuration.Configuration{
		Output: "file",
		Install: configuration.Install{
			Path: "/root",
		},
	}
	err := validateConfiguration(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrImageRequired, err)

}

func TestValidateConfigurationEmptyPath(t *testing.T) {

	c := &configuration.Configuration{
		Output: "file",
		Install: configuration.Install{
			Image: "debian",
		},
	}
	err := validateConfiguration(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrPathRequired, err)

}

type DockerMock struct {
	mock.Mock
}

func (m *DockerMock) Run(option docker.RunOptions) (int, error) {
	args := m.Called(option)
	return args.Int(0), args.Error(1)
}

func (m *DockerMock) Ping() error {
	return nil
}

func (m *DockerMock) Logs(id string, stream docker.LogStream) error {
	return nil
}

func (m *DockerMock) Build(option docker.BuildOptions) error {
	return nil
}

func TestInstall(t *testing.T) {

	c, o := dockerMockConf("../testing/app.bin")

	d := &DockerMock{}
	d.On("Run", o).Return(0, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.Nil(t, err)

}

func TestInstallOnError(t *testing.T) {

	c, o := dockerMockConf("../testing/app.bin")

	d := &DockerMock{}
	d.On("Run", o).Return(255, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestInstallWithDockerError(t *testing.T) {

	c, o := dockerMockConf("../testing/app.bin")

	d := &DockerMock{}
	e := errors.New("an error")
	d.On("Run", o).Return(0, e)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, e, err)

}

func TestInstallWithConfigurationError(t *testing.T) {

	c := &configuration.Configuration{}
	d := &DockerMock{}

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func TestInstallWithNilConfiguration(t *testing.T) {

	d := &DockerMock{}

	err := Install(d, nil)

	d.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, ErrConfigurationEmpty, err)

}

func TestInstallWithNoOutput(t *testing.T) {

	c, o := dockerMockConf("file.txt")

	d := &DockerMock{}
	d.On("Run", o).Return(0, nil)

	err := Install(d, c)

	d.AssertExpectations(t)
	assert.NotNil(t, err)

}

func dockerMockConf(output string) (*configuration.Configuration, docker.RunOptions) {

	c := &configuration.Configuration{
		Output: output,
		Install: configuration.Install{
			Command: "make",
			Path:    "/root",
			Image:   "debian",
		},
	}

	e, _ := GetEnv(c)
	v, _ := GetVolumes(c)

	o := docker.RunOptions{
		Image:   "debian",
		Command: "make",
		Env:     e,
		Volumes: v,
	}

	return c, o
}
