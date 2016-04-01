package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type streamDecoderWrapper struct {
	out    *bufio.Writer
	err    *bufio.Writer
	buffer *bytes.Buffer
	tty    bool
	lines  map[string]int
	diff   int
}

type message struct {
	Stream   string `json:"stream,omitempty"`
	Status   string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	ID       string `json:"id,omitempty"`
	Error    string `json:"error,omitempty"`
}

func newStreamDecoderWrapper(out io.Writer, err io.Writer, tty bool) io.Writer {

	buffer := &bytes.Buffer{}
	lines := make(map[string]int)
	diff := 0

	return &streamDecoderWrapper{
		out: bufio.NewWriter(out), err: bufio.NewWriter(err), buffer: buffer,
		tty: tty, lines: lines, diff: diff,
	}
}

func (w *streamDecoderWrapper) Write(p []byte) (int, error) {

	defer w.Flush()

	n, err := w.buffer.Write(p)
	if err != nil {
		return n, err
	}

	return n, w.Decode()
}

func (w streamDecoderWrapper) Flush() {
	_ = w.err.Flush()
	_ = w.out.Flush()
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

func (w *streamDecoderWrapper) setCursor(id string) {

	if id != "" {

		line, ok := w.lines[id]
		if !ok {

			// NOTE: This approach of using len(lines) to
			// figure out the number of lines of history
			// only works as long as we clear the history
			// when we output something that's not
			// accounted for in the map, such as a line
			// with no ID.

			line = len(w.lines)
			w.lines[id] = line

			if w.tty {
				fmt.Fprintf(w.getOutWriter(), "\n")
			}

		} else {
			w.diff = len(w.lines) - line
		}

		if w.tty {

			// NOTE: this appears to be necessary even if
			// diff == 0.
			// <ESC>[{diff}A = move cursor up diff rows
			fmt.Fprintf(w.getOutWriter(), "%c[%dA", 27, w.diff)

		}

	} else {

		// When outputting something that isn't progress
		// output, clear the history of previous lines. We
		// don't want progress entries from some previous
		// operation to be updated (for example, pull -a
		// with multiple tags).

		w.lines = make(map[string]int)

	}
}

func (w *streamDecoderWrapper) Decode() error {

	decoder := json.NewDecoder(w.buffer)

	for {

		m := message{}

		if err := decoder.Decode(&m); err != nil {
			if err != io.EOF {
				return err
			}

			// Recopy remaining bytes into buffer to be available again for json decoder.
			b, err := ioutil.ReadAll(decoder.Buffered())
			if err != nil {
				return err
			}

			w.buffer = bytes.NewBuffer(b)
			return nil
		}

		if m.Error != "" {
			return errors.New(m.Error)
		}

		w.setCursor(m.ID)
		w.encode(m)

		if m.ID != "" && w.tty {
			// NOTE: this appears to be necessary even if
			// diff == 0.
			// <ESC>[{diff}B = move cursor down diff rows
			fmt.Fprintf(w.getOutWriter(), "%c[%dB", 27, w.diff)
		}

	}

}
