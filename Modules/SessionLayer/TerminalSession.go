package SessionLayer

import (
	"fmt"
	"path/filepath"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/API"
	"sssh_server/Modules/SSH"
	"time"
)

/*
safe for copy
underlying types use pointers to data
*/
type TerminalSession struct {
	ID            string
	SSHSessions   map[string]*SSH.SSHSession    `json:"-"`
	HandlersMap   map[string]API.SessionHandler `json:"-"`
	Modules       []API.Module                  `json:"-"`
	Config        API.SessionConfig             `json:"-"`
	TimeCreated   int64
	LastConnected int64
	Name          string
}

// Create a terminal TerminalSession from a ssh session
func newSession(id string) (s *TerminalSession) {
	s = new(TerminalSession)
	s.SSHSessions = make(map[string]*SSH.SSHSession)
	s.HandlersMap = make(map[string]API.SessionHandler)
	s.ID = id
	s.Name = id
	//s.recentCommandsMutex.Lock()
	//s.InitTerminal()
	s.Modules = []API.Module{}
	s.InitConfig()
	s.TimeCreated = time.Now().Unix()
	s.LastConnected = time.Now().Unix()

	addServices(s)

	return s
}

func (s *TerminalSession) GetID() string {
	return s.ID
}

func (s *TerminalSession) GetConfig() API.SessionConfig {
	return s.Config
}

func (s *TerminalSession) InitConfig() {
	s.Config = API.SessionConfig{}
	s.Config.SessionID = s.ID

	basePath, err := filepath.Abs("Assets/")
	CustomUtils.CheckPanic(err, "Could not create history file for the session")
	historyFilePath := fmt.Sprintf("%v/%v", basePath, HISTORY_FILE_NAME+s.ID)
	s.Config.HistoryFilePath = historyFilePath
}

func (s *TerminalSession) OnNewSessionLifecycleHook() {
	// Call lifecycle hooks on the service
	for _, service := range s.Modules {
		service.OnNewSession(s)
	}
}

func (s *TerminalSession) OnNewConnectionLifecycleHook(sshSession *SSH.SSHSession) {
	// Call lifecycle hooks on the service
	for _, service := range s.Modules {
		service.OnNewConnection(sshSession)
	}
}

func addServices(terminalSession *TerminalSession) {
	for _, service := range Services {
		terminalSession.addService(service)
	}
}

// Internal function to add services
func (s *TerminalSession) addService(service API.Module) {
	s.Modules = append(s.Modules, service)
	handlers := service.GetHandlers()
	for _, handler := range handlers {
		s.HandlersMap[handler.Name] = handler.RequestHandler
	}
}

var _ API.TerminalSessionInterface = (*TerminalSession)(nil)
