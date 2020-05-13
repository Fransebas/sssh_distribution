package CommandExecuter

import (
	"fmt"
	"github.com/kr/pty"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/SessionModules/API"
)

func (cme *CommandExecuter) GetHandlers() []*API.RequestHandler {
	// Executes a command using the ssh connection as god intended to be
	execCommandOnce := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, err := CustomUtils.Read(r)
			CustomUtils.CheckPrint(err)
			//fmt.Println("exec command: " + string(b))
			_, _ = w.Write([]byte(CustomUtils.ExecuteCommandOnce(string(b), cme.user)))
		},
		Name: "exec.once",
	}

	// Executes a command using the ssh connection as god intended to be
	execCommandOpen := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, err := CustomUtils.Read(r)

			CustomUtils.CheckPrint(err)
			//fmt.Println("exec command: " + string(b))

			c := CustomUtils.ExecuteCommand(string(b), cme.user)

			pt, e := pty.Start(c)

			CustomUtils.CheckPrint(e)

			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered in f", r)
					}
				}()
				if c.Stdin != nil {
					_, _ = io.Copy(w, pt)
				}
			}()
			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered in f", r)
					}
				}()
				if c.Stdout != nil {
					_, _ = io.Copy(pt, r)
				}
			}()
		},
		Name: "exec.on",
	}

	// Executes a command using the ssh connection as god intended to be
	manCommand := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, err := CustomUtils.Read(r)
			CustomUtils.CheckPrint(err)
			//fmt.Println("exec command: " + string(b))
			manCmnd := fmt.Sprintf("man %s | col -b", string(b))
			_, _ = w.Write([]byte(cme.ExecuteCommand(manCmnd)))
		},
		Name: "man",
	}

	return []*API.RequestHandler{&execCommandOnce, &manCommand, &execCommandOpen}
}
