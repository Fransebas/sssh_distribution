package TerminalService

import (
	"sssh_server/Services/API"
	"sssh_server/Services/SSH"
	"sssh_server/Terminal"
)

type TerminalService struct {
	Terminal *Terminal.Terminal
}

func (ts *TerminalService) OnNewSession(session API.TerminalSessionInterface) {
	config := session.GetConfig()
	ts.Terminal = Terminal.InitTerminal(session.GetID(), config.GetHistoryFilePath())
	ts.Terminal.Run()
}

func (ts *TerminalService) OnNewConnection(sshSession *SSH.SSHSession) {

}

var _ API.Service = (*TerminalService)(nil) // Verify that *T implements I.
