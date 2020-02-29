package Terminal

import (
	"sssh_server/CustomUtils"
)

// The task of the TerminalReader is to encapsulate the read of the terminal
// such that multiple clients can read the input simultaneously
// Basically every client should have a TerminalReader which will read from a shared buffer
// and the TerminalReader remembers where each client is currently reading from
type TerminalReader struct {
	buffer   *CustomUtils.FixedDeque
	offset   int // current reading point
	terminal *Terminal
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Read that makes TerminalReader implements the interface of io.Reader
func (tr *TerminalReader) Read(p []byte) (n int, err error) {
	// Update the buffer if necesary
	tr.terminal.read()

	n, err = tr.BufferRead(p)
	//CustomUtils.LogTime(string(p[:n]), "buffer read / send")
	return
}

// Read the following bytes, from tr.offset to end and update offset
func (tr *TerminalReader) BufferRead(p []byte) (n int, err error) {

	var b []byte
	b, tr.offset = tr.buffer.BytesFrom(tr.offset)
	n = min(len(p), len(b))
	for i := 0; i < n; i++ {
		p[i] = b[i]
	}
	return n, nil
}
