package API

import (
	"sssh_server/Services/SSH"
	"sssh_server/Services/SessionLayer"
)

type Service interface {
	OnNewSession(s *SessionLayer.TerminalSession)
	OnNewConnection(sshSession *SSH.SSHSession)
	GetHandlers() []SessionLayer.SessionHandler
	//SetHandlers(sessionService *SessionLayer.SessionService)
}
