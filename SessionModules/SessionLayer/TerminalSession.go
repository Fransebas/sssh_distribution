package SessionLayer

import (
	"sssh_server/Modules/DirectoryManager"
	"sssh_server/Modules/SSH"
	"sssh_server/SessionModules/API"
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
	Username      string
}

// Create a terminal TerminalSession from a ssh session
func newSession(id, username string) (s *TerminalSession) {
	s = new(TerminalSession)
	s.SSHSessions = make(map[string]*SSH.SSHSession)
	s.HandlersMap = make(map[string]API.SessionHandler)
	s.ID = id
	s.Name = id
	//s.recentCommandsMutex.Lock()
	//s.InitTerminal()
	s.Modules = []API.Module{}
	s.TimeCreated = time.Now().Unix()
	s.LastConnected = time.Now().Unix()
	s.Username = username
	s.InitConfig()

	addServices(s)

	return s
}

func (s *TerminalSession) GetUsername() string {
	return s.Username
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

	dm := DirectoryManager.New(s.Username)

	s.Config.HistoryFilePath = dm.GetVariableFile(HISTORY_FILE_NAME + s.ID)
	s.Config.BashrcFilePath = dm.GetConfigFile("bashrc")
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

func (s *TerminalSession) Close() {
	// Call lifecycle hooks on the service
	for _, module := range s.Modules {
		module.Close()
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
