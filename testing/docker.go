package testing

import (
	"github.com/crowley-io/pack/docker"
	"github.com/stretchr/testify/mock"
)

// DockerMock is a mock of docker.Docker interface.
type DockerMock struct {
	mock.Mock
}

// Run is a mock function for docker.Docker interface.
func (m *DockerMock) Run(option docker.RunOptions, stream docker.LogStream) (int, error) {
	args := m.Called(option, stream)
	return args.Int(0), args.Error(1)
}

// Ping is a mock function for docker.Docker interface.
func (m *DockerMock) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// Logs is a mock function for docker.Docker interface.
func (m *DockerMock) Logs(id string, stream docker.LogStream) error {
	args := m.Called(id, stream)
	return args.Error(0)
}

// Build is a mock function for docker.Docker interface.
func (m *DockerMock) Build(option docker.BuildOptions, stream docker.LogStream) error {
	args := m.Called(option, stream)
	return args.Error(0)
}

// Tag is a mock function for docker.Docker interface.
func (m *DockerMock) Tag(option docker.TagOptions) error {
	args := m.Called(option)
	return args.Error(0)
}

// Push is a mock function for docker.Docker interface.
func (m *DockerMock) Push(option docker.PushOptions, stream docker.LogStream) error {
	args := m.Called(option, stream)
	return args.Error(0)
}

// ImageID is a mock function for docker.Docker interface.
func (m *DockerMock) ImageID(name string) string {
	args := m.Called(name)
	return args.String(0)
}

// RemoveImage is a mock function for docker.Docker interface.
func (m *DockerMock) RemoveImage(name string) error {
	args := m.Called(name)
	return args.Error(0)
}
