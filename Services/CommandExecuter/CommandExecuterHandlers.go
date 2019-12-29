package CommandExecuter

import (
	"fmt"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
)

func (cme *CommandExecuter) GetHandlers() []*API.RequestHandler {
	// Executes a command using the ssh connection as god intended to be
	execCommand := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, err := CustomUtils.Read(r)
			CustomUtils.CheckPrint(err)
			//fmt.Println("exec command: " + string(b))
			_, _ = w.Write([]byte(cme.ExecuteCommand(string(b))))
		},
		Name: "exec",
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

	return []*API.RequestHandler{&execCommand, &manCommand}
}
