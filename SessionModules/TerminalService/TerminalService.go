package TerminalService

import (
	"sssh_server/Modules/SSH"
	"sssh_server/Modules/Terminal"
	"sssh_server/SessionModules/API"
)

type TerminalService struct {
	Terminal *Terminal.Terminal
}

func (ts *TerminalService) OnNewSession(session API.TerminalSessionInterface) {
	config := session.GetConfig()
	ts.Terminal = Terminal.InitTerminal(session.GetID(), config.GetHistoryFilePath(), session.GetUsername())
	ts.Terminal.Run()
}

func (ts *TerminalService) OnNewConnection(sshSession *SSH.SSHSession) {

}

var _ API.Module = (*TerminalService)(nil) // Verify that *T implements I.
