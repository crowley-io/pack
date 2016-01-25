package docker

import (
	"encoding/json"
	"fmt"
	"io"
)

type streamDecoderWrapper struct {
	out     io.Writer
	err     io.Writer
	decoder *json.Decoder
	reader  *io.PipeReader
	writer  *io.PipeWriter
	tty     bool
}

type message struct {
	Stream   string `json:"stream,omitempty"`
	Status   string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	ID       string `json:"id,omitempty"`
}

func newStreamDecoderWrapper(out io.Writer, err io.Writer) io.WriteCloser {

	reader, writer := io.Pipe()
	decoder := json.NewDecoder(reader)

	wrapper := &streamDecoderWrapper{
		out: out, err: err, decoder: decoder,
		reader: reader, writer: writer, tty: true,
	}

	go wrapper.Flush()
	return wrapper

}

func (w streamDecoderWrapper) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w streamDecoderWrapper) Close() error {
	return w.writer.Close()
}

func (w streamDecoderWrapper) getOutWriter() io.Writer {
	return w.out
}

func (w streamDecoderWrapper) getErrWriter() io.Writer {
	return w.err
}

func (w streamDecoderWrapper) encode(m message) {

	var endl string
	if w.tty && m.Stream == "" && (m.Progress != "" || m.Status != "") {
		// <ESC>[2K = erase entire current line
		fmt.Fprintf(w.getOutWriter(), "%c[2K\r", 27)
		endl = "\r"
	} else if m.Progress != "" { //disable progressbar in non-terminal
		return
	}

	if m.ID != "" {
		fmt.Fprintf(w.getOutWriter(), "%s: ", m.ID)
	}
	if m.Progress != "" && w.tty {
		fmt.Fprintf(w.getOutWriter(), "%s %s%s", m.Status, m.Progress, endl)
	} else if m.Stream != "" {
		fmt.Fprintf(w.getOutWriter(), "%s%s", m.Stream, endl)
	} else {
		fmt.Fprintf(w.getOutWriter(), "%s%s\n", m.Status, endl)
	}
}

func (w streamDecoderWrapper) setCursor(id string, lines *(map[string]int), diff *(int)) {

	if id != "" {

		line, ok := (*lines)[id]
		if !ok {

			// NOTE: This approach of using len(lines) to
			// figure out the number of lines of history
			// only works as long as we clear the history
			// when we output something that's not
			// accounted for in the map, such as a line
			// with no ID.

			line = len(*lines)
			(*lines)[id] = line

			if w.tty {
				fmt.Fprintf(w.getOutWriter(), "\n")
			}

		} else {
			(*diff) = len(*lines) - line
		}

		if w.tty {

			// NOTE: this appears to be necessary even if
			// diff == 0.
			// <ESC>[{diff}A = move cursor up diff rows

			fmt.Fprintf(w.getOutWriter(), "%c[%dA", 27, (*diff))

		}

	} else {

		// When outputting something that isn't progress
		// output, clear the history of previous lines. We
		// don't want progress entries from some previous
		// operation to be updated (for example, pull -a
		// with multiple tags).

		(*lines) = make(map[string]int)

	}
}

func (w streamDecoderWrapper) Flush() {

	lines := make(map[string]int)

	for {

		m := message{}
		diff := 0

		if err := w.decoder.Decode(&m); err != nil {
			if err != io.EOF {
				fmt.Fprintln(w.getErrWriter(), err)
			}
			if err := w.reader.Close(); err != nil {
				fmt.Fprintln(w.getErrWriter(), err)
			}
			return
		}

		w.setCursor(m.ID, &lines, &diff)
		w.encode(m)

		if m.ID != "" && w.tty {
			// NOTE: this appears to be necessary even if
			// diff == 0.
			// <ESC>[{diff}B = move cursor down diff rows
			fmt.Fprintf(w.getOutWriter(), "%c[%dB", 27, diff)
		}

	}

}
