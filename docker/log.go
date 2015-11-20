package docker

import (
	"io"

	api "github.com/fsouza/go-dockerclient"
)

// LogStream contains two io.Writer for respectively, stdout and stderr.
type LogStream struct {
	Out io.Writer
	Err io.Writer
}

// See Docker interface
func (d docker) Logs(id string, stream LogStream) error {
	return d.client.Logs(logsOptions(id, stream))
}

func logsOptions(container string, stream LogStream) api.LogsOptions {
	return api.LogsOptions{
		Container:    container,
		OutputStream: stream.Out,
		ErrorStream:  stream.Err,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   false,
	}
}
