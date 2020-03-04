package TerminalService

import (
	"fmt"
	"sssh_server/Modules/SSH"
	"sssh_server/Modules/Terminal"
	"sssh_server/SessionModules/API"
)

type TerminalService struct {
	Terminal  *Terminal.Terminal
	pwd       string
	pwdChange func(path string)
}

func (ts *TerminalService) OnNewSession(session API.TerminalSessionInterface) {
	config := session.GetConfig()
	ts.Terminal = Terminal.InitTerminal(session.GetID(), config.GetHistoryFilePath(), session.GetUsername())
	ts.Terminal.Run()

}

func (ts *TerminalService) OnNewConnection(sshSession *SSH.SSHSession) {

}

func (ts *TerminalService) Close() {

}

func (ts *TerminalService) OnPWD(path string) {
	ts.pwd = path
	fmt.Printf("\nPWD %v\n", path)

	if ts.pwdChange != nil {
		ts.pwdChange(ts.pwd)
	}
}

var _ API.Module = (*TerminalService)(nil) // Verify that *T implements I.
