package SessionLayer

import (
	"fmt"
	"io"
	"sssh_server/CustomUtils"
	"sssh_server/Services/CommandList"
	"sssh_server/Services/GlobalVariables"
	"sssh_server/Services/RecentCommands"
	"sssh_server/Services/SSH"
	"sssh_server/Terminal"
	"sync"
)

/*
safe for copy
underlying types use pointers to data
*/
type TerminalSession struct {
	ID                  string
	Terminal            *Terminal.Terminal
	recentCommands      *RecentCommands.RecentCommands
	recentCommandsMutex sync.Mutex
	commandList         *CommandList.CommandList
	SSHSessions         map[string]*SSH.SSHSession
	GlobalVars          GlobalVariables.GlobalVariables
}

type SessionService struct {
	Server            *SSH.SSSHServer
	Sessions          map[string]*TerminalSession
	SSHidToTerminalID map[string]string
}

type SessionHandler func(s *TerminalSession, w io.Writer, r io.Reader)

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s *SessionService) {
	s = new(SessionService)

	s.Server = &SSH.SSSHServer{}
	s.Sessions = make(map[string]*TerminalSession)
	s.SSHidToTerminalID = make(map[string]string)

	s.InitHandlers()

	s.Server.OnNewSession(func(anonymousSession *SSH.SSHSession) {
		// When a new SSH session is created, then create a terminal session

		fmt.Println("Init ssh")
		s.Server.HandleFunc("session", func(sshSession *SSH.SSHSession, w io.Writer, r io.Reader) {

			// TODO: handle the null terminalSession number as a new terminalSession
			b, _ := CustomUtils.Read(r)
			sessionID := string(b)

			fmt.Println("Terminal session started " + sessionID)

			if s.Sessions[sessionID] != nil {
				terminalSession := s.Sessions[sessionID]
				terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession
			} else {
				terminalSession := newSession(sshSession, sessionID)
				s.Sessions[sessionID] = terminalSession
				terminalSession.SSHSessions[sshSession.GetSessionID()] = sshSession
			}
			s.SSHidToTerminalID[sshSession.GetSessionID()] = sessionID

			w.Write([]byte("\n"))
		})
	})

	return s
}

func (s *SessionService) SSHSessionToTerminalSession(sshSession *SSH.SSHSession) *TerminalSession {
	return s.Sessions[s.SSHidToTerminalID[sshSession.GetSessionID()]]
}

func (s *SessionService) HandleFunc(msgType string, handler SessionHandler) {
	s.Server.HandleFunc(msgType, func(session *SSH.SSHSession, w io.Writer, r io.Reader) {
		handler(s.SSHSessionToTerminalSession(session), w, r)
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

		fmt.Println("2")
		//session.recentCommandsMutex.Lock()
		newCommands := session.recentCommands.UpdateRecentCommands(data)
		fmt.Printf("Number of commands= %v \n", len(newCommands))
		//session.recentCommandsMutex.Unlock()

		//newCommandsJson, err := json.Marshal(newCommands)
		//if err == nil {
		//
		//	socket.SocketConn.Emit("commands", string(newCommandsJson))
		//}
	}
}

// Create a terminal TerminalSession from a ssh session
func newSession(sshSession *SSH.SSHSession, id string) (s *TerminalSession) {
	s = new(TerminalSession)
	s.SSHSessions = make(map[string]*SSH.SSHSession)
	s.ID = id
	//s.recentCommandsMutex.Lock()
	s.InitTerminal()
	s.commandList = CommandList.NewCommandList()
	s.GlobalVars = GlobalVariables.GlobalVariables{}
	return s
}

func (s *TerminalSession) InitTerminal() {
	fmt.Printf("ID here = %v\n", s.ID)
	s.Terminal = Terminal.InitTerminal(s.ID)
	s.Terminal.Run()
}

func (s *TerminalSession) InitCommand() {

}
