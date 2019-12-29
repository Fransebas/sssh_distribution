package API

import (
	"io"
	"sssh_server/Services/SSH"
)

type SessionHandler func(w io.Writer, r io.Reader)

type TerminalSessionInterface interface {
	GetID() string
	GetConfig() SessionConfig
}

type Service interface {
	OnNewSession(s TerminalSessionInterface)
	OnNewConnection(sshSession *SSH.SSHSession)
	GetHandlers() []*RequestHandler
	//GetHTTPHandlers()
	//SetHandlers(sessionService *SessionLayer.SessionService)
}

type HistoryService interface {
	OnNewCommand(cmnd string)
}

type RequestHandler struct {
	RequestHandler SessionHandler
	Name           string
}

type HTTPHandler struct {
	RequestHandler SessionHandler
	Name           string
}
