package TerminalService

import (
	"encoding/json"
	"fmt"
	"github.com/creack/pty"

	"io"
	"sssh_server/CustomUtils"
	"sssh_server/SessionModules/API"
)

var c = 0

func (ts *TerminalService) GetHandlers() []*API.RequestHandler {
	terminalConnection := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			c++
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
				_, e := io.Copy(ts.Terminal, r)
				fmt.Printf("Writing Connection Lost %v", e)
			}()

			b := make([]byte, 1024*8)
			n, _ := reader.BufferRead(b)
			_, e := w.Write(b[:n])

			CustomUtils.CheckPrint(e)

			//buf := make([]byte, 3*1024)
			//for {
			//	nr, e := reader.Read(buf)
			//	fmt.Printf("w %v \n", c)
			//	if e != nil {
			//		fmt.Println("Connection Lost")
			//		break;
			//	}
			//	if nr > 0 {
			//		_, _ = w.Write(buf[0:nr])
			//	}
			//}

			_, e = io.Copy(w, &reader)
			fmt.Printf("Reading Connection Lost %v", e)
		},
		Name: "terminal",
	}

	terminalResize := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// On terminal connection opened

			var resize pty.Winsize
			b, err := CustomUtils.Read(r)
			err = json.Unmarshal(b, &resize)

			CustomUtils.CheckPrint(err)

			if err == nil {
				ts.Terminal.SetSize(&resize)
			}
		},
		Name: "terminal.resize",
	}

	terminalPWD := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			// Send the current working directory every time it changes
			_, e := w.Write([]byte(ts.pwd))
			CustomUtils.CheckPrint(e)

			ts.pwdChange = func(path string) {
				_, _ = w.Write([]byte(ts.pwd))
			}

		},
		Name: "terminal.pwd",
	}

	return []*API.RequestHandler{&terminalConnection, &terminalResize, &terminalPWD}
}
