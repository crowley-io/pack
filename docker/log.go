package docker

import (
	"io"
	"os"

	"github.com/andrew-d/go-termutil"
	api "github.com/fsouza/go-dockerclient"
)

// LogStream contains two io.Writer for respectively, stdout and stderr, and also a JSON Decoder for the Docker API.
type LogStream struct {
	Out     io.Writer
	Err     io.Writer
	Decoder io.Writer
}

// OutputStream return a stream output handler with its encoding.
func (l LogStream) OutputStream() (json bool, stream io.Writer) {
	return true, l.Decoder
}

// See Docker interface
func (d docker) Logs(id string, stream LogStream) error {
	return d.client.Logs(logsOptions(id, stream))
}

// NewLogStream return a default LogStream using OS stdout and stderr.
func NewLogStream() LogStream {
	out := os.Stdout
	err := os.Stderr
	decoder := newStreamDecoderWrapper(out, err, termutil.Isatty(out.Fd()))
	return LogStream{Out: out, Err: err, Decoder: decoder}
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
