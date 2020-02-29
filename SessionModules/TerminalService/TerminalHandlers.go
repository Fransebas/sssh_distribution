package TerminalService

import (
	"encoding/json"
	"fmt"
	"github.com/creack/pty"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/SessionModules/API"
	"time"
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
			go func() {
				_, _ = io.Copy(ts.Terminal, r)
			}()

			b := make([]byte, 1024*8)
			n, _ := reader.BufferRead(b)
			_, e := w.Write(b[:n])

			CustomUtils.CheckPrint(e)

			buf := make([]byte, 3*1024)
			for {
				CustomUtils.LogTime("a", "go.before.r")
				nr, _ := reader.Read(buf)
				time.Sleep(100 * time.Millisecond)
				CustomUtils.LogTime("a", "go.after.r")
				if nr > 0 {
					fmt.Printf("data %v\n", string(buf[0:nr]))
					_, _ = w.Write(buf[0:nr])
					CustomUtils.LogTime("a", "go.after.w")
				}
			}

			//_, _ = io.Copy(w, &reader)
		},
		Name: "terminal",
	}

	terminalResize := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// On terminal connection opened
			var resize pty.Winsize
			b, err := CustomUtils.Read(r)
			err = json.Unmarshal(b, &resize)
			if err != nil {
				ts.Terminal.SetSize(&resize)
			}
		},
		Name: "terminal.resize",
	}

	return []*API.RequestHandler{&terminalConnection, &terminalResize}
}
