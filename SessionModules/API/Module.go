package API

import (
	"io"
	"sssh_server/Modules/SSH"
)

type SessionHandler func(w io.Writer, r io.Reader)

type TerminalSessionInterface interface {
	GetID() string
	GetConfig() SessionConfig
	GetUsername() string
}

type Module interface {
	OnNewSession(s TerminalSessionInterface)
	OnNewConnection(sshSession *SSH.SSHSession)
	GetHandlers() []*RequestHandler
	Close()
	//GetHTTPHandlers()
	//SetHandlers(sessionService *SessionLayer.SessionService)
}

// This function should run on the client
type ClientCode func(cmnd string) string

//type ClientCode interface {
//	OnCommandRun(cmnd string) string
//}

// TODO: not working
// this interface is if your service needs specific information from "inside" the bash session,
// for example, the environment variables or the last run command
type OnCommandService interface {
	// Get the function to run on the client
	GetClientCode() ClientCode
	GetName() string
	// Receives the resulting string from the ClientCode
	ClientResponse(res string)
}

type HistoryService interface {
	OnNewCommand(cmnd string)
}

// Runs every time the user types a command and returns the current working directory
type PWDService interface {
	OnPWD(path string)
}

type VariablesService interface {
	OnUpdateVariables(vars string)
}

type RequestHandler struct {
	RequestHandler SessionHandler
	Name           string
}

type HTTPHandler struct {
	RequestHandler SessionHandler
	Name           string
}
