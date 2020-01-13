package EchoService

import (
	"fmt"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/API"
	"sssh_server/Modules/SSH"
)

type EchoService struct {
}

func (e *EchoService) OnNewSession(session API.TerminalSessionInterface) {

}
func (e *EchoService) OnNewConnection(sshSession *SSH.SSHSession) {

}
func (e *EchoService) GetHandlers() []*API.RequestHandler {
	echoCommand := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			b, _ := CustomUtils.Read(r)
			fmt.Println("echoing : " + string(b))
			_, e := w.Write(b)
			CustomUtils.CheckPrint(e)
		},
		Name: "echo",
	}

	echoConnection := API.RequestHandler{
		RequestHandler: func(w io.Writer, r io.Reader) {
			_, _ = io.Copy(w, r)
		},
		Name: "echoConnection",
	}

	return []*API.RequestHandler{&echoCommand, &echoConnection}
}
