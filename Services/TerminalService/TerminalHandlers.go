package TerminalService

import (
	"encoding/json"
	"github.com/creack/pty"
	"io"
	"io/ioutil"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
)

func (ts *TerminalService) GetHandlers() []*API.RequestHandler {
	terminalConnection := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// On terminal connection opened
			// Simple redirect the channels

			// w is the user terminal
			// r is the user input on her terminal

			//w.Write(ts.Terminal.GetBuffer())
			reader := ts.Terminal.GetReader()

			// TODO: make sure we don't have a leak here
			// the documentation here https://golang.org/pkg/io/#Copy
			// says that Copy stops when the src returns EOF but not when the
			// writer hence this should stop in only half of the cases
			go func() { _, _ = io.Copy(ts.Terminal, r) }()

			b := make([]byte, 1024*8)
			n, _ := reader.BufferRead(b)
			_, e := w.Write(b[:n])

			CustomUtils.CheckPrint(e)

			_, _ = io.Copy(w, &reader)
		},
		Name: "terminal",
	}

	terminalResize := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// On terminal connection opened
			var resize pty.Winsize
			b, err := ioutil.ReadAll(r)
			err = json.Unmarshal(b, &resize)
			if err != nil {
				ts.Terminal.SetSize(&resize)
			}
		},
		Name: "terminal.resize",
	}

	return []*API.RequestHandler{&terminalConnection, &terminalResize}
}
