package SessionLayer

import (
	"fmt"
	"io"
	"path/filepath"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/CommandList"
	"sssh_server/Services/EchoService"
	"sssh_server/Services/GlobalVariables"
	"sssh_server/Services/History"
	"sssh_server/Services/SSH"
	"sssh_server/Services/TerminalService"
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
	Services      []API.Service                 `json:"-"`
	Config        API.SessionConfig             `json:"-"`
	TimeCreated   int64
	LastConnected int64
	Name          string
}

type SessionService struct {
	Server            *SSH.SSSHServer
	Sessions          map[string]*TerminalSession
	SSHidToTerminalID map[string]string
}

var HISTORY_FILE_NAME = "history"

func Constructor() (s *SessionService) {
	s = new(SessionService)

	s.Server = &SSH.SSSHServer{}
	s.Sessions = make(map[string]*TerminalSession)
	s.SSHidToTerminalID = make(map[string]string)

	s.ChannelHandler()

	s.Server.OnNewSession(func(anonymousSession *SSH.SSHSession) {
		//
	})

	return s
}

func (s *SessionService) SSHSessionToTerminalSession(sshSession *SSH.SSHSession) *TerminalSession {
	return s.Sessions[s.SSHidToTerminalID[sshSession.GetSessionID()]]
}

func (s *SessionService) CreateSession(msgType string, sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {
	// TODO: handle the null terminalSession number as a new terminalSession
	b, _ := CustomUtils.Read(r)
	sessionID := string(b)

	fmt.Println("Terminal session started " + sessionID)

	var terminalSession *TerminalSession

	if s.Sessions[sessionID] != nil {
		// Terminal Session already exists, add sshSession to it
		terminalSession = s.Sessions[sessionID]
		terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession
	} else {
		// Terminal Session doesn't exist, create one
		terminalSession = newSession(sessionID)
		s.Sessions[sessionID] = terminalSession
		terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession

		// Lifecycle hook
		terminalSession.OnNewSessionLifecycleHook()
	}
	terminalSession.LastConnected = time.Now().Unix()

	// Lifecycle hook
	terminalSession.OnNewConnectionLifecycleHook(sshSession)

	s.SSHidToTerminalID[sshSession.GetSessionID()] = sessionID

	w.Write([]byte("\n"))
}

// Make the mapping for any kind of messages
// Here is where the multiplexing of channels happen
func (s *SessionService) ChannelHandler() {
	s.Server.SetAnyHandler(func(msgType string, sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {

		if msgType == "session" {
			// On the session message create a new session
			s.CreateSession(msgType, sshSession, w, r)
		} else {
			terminalSession := s.SSHSessionToTerminalSession(sshSession)
			if handler, ok := terminalSession.HandlersMap[msgType]; ok {
				handler(w, r)
			} else {
				CustomUtils.CheckPrint(fmt.Errorf("Message type " + msgType + " doesn't exist"))
			}
		}
	})
}

func (s *SessionService) Serve() {
	s.Server.ListenAndServe()
}

// Adds a new command to the recently used commands
func (ss *SessionService) AddCommand(data string, id string) {
	fmt.Printf("data %v \n", data)
	if _, ok := ss.Sessions[id]; ok {
		session := ss.Sessions[id]
		for _, service := range session.Services {
			if historyService, ok := service.(API.HistoryService); ok {
				historyService.OnNewCommand(data)
			}
		}
	}
}

func (s *TerminalSession) addService(service API.Service) {
	s.Services = append(s.Services, service)
	handlers := service.GetHandlers()
	for _, handler := range handlers {
		s.HandlersMap[handler.Name] = handler.RequestHandler
	}
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
	s.Services = []API.Service{}
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
	for _, service := range s.Services {
		service.OnNewSession(s)
	}
}

func (s *TerminalSession) OnNewConnectionLifecycleHook(sshSession *SSH.SSHSession) {
	// Call lifecycle hooks on the service
	for _, service := range s.Services {
		service.OnNewConnection(sshSession)
	}
}

func addServices(terminalSession *TerminalSession) {
	terminalSession.addService(new(TerminalService.TerminalService))
	terminalSession.addService(new(History.History))
	terminalSession.addService(new(CommandList.CommandListService))
	terminalSession.addService(new(CommandExecuter.CommandExecuter))
	terminalSession.addService(new(EchoService.EchoService))
	terminalSession.addService(new(GlobalVariables.GlobalVariables))
}

var _ API.TerminalSessionInterface = (*TerminalSession)(nil)
