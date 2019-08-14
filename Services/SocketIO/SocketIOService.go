package SocketIO

import (
	"fmt"
	"rgui/Services/CommandList"
	"rgui/Services/RecentCommands"
	"rgui/Services/SSH"
	"rgui/Terminal"
	"sync"
)

/*
safe for copy
underlying types use pointers to data
*/
type Session struct {
	ID                  string
	Terminal            *Terminal.Terminal
	recentCommands      *RecentCommands.RecentCommands
	recentCommandsMutex sync.Mutex
	commandList         *CommandList.CommandList
	SSHSession          *SSH.Session
}

type SocketIOService struct {
	Server   *SSH.SSSHServer
	Sessions map[string]*Session
}

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s *SocketIOService) {
	s = new(SocketIOService)

	s.Server = &SSH.SSSHServer{}
	s.Sessions = make(map[string]*Session)

	s.InitHandlers()

	s.Server.OnNewSession(func(shhSession *SSH.Session) {
		session := newSession(shhSession)
		s.Sessions[shhSession.GetSessionID()] = session
	})

	return s
}

func (s *SocketIOService) Serve() {
	s.Server.ListenAndServe()
}

// Adds a new command to the recently used commands
func (ss *SocketIOService) AddCommand(data string, id string) {

	fmt.Printf("data %v \n", data)
	if _, ok := ss.Sessions[id]; ok {
		session := ss.Sessions[id]

		fmt.Println("2")
		session.recentCommandsMutex.Lock()
		newCommands := session.recentCommands.UpdateRecentCommands(data)
		fmt.Printf("Number of commands= %v \n", len(newCommands))
		session.recentCommandsMutex.Unlock()
		//newCommandsJson, err := json.Marshal(newCommands)
		//if err == nil {
		//
		//	socket.SocketConn.Emit("commands", string(newCommandsJson))
		//}
	}
}

func newSession(shhSession *SSH.Session) (s *Session) {
	s = new(Session)
	s.ID = shhSession.GetSessionID()
	s.SSHSession = shhSession
	fmt.Println("1")
	s.recentCommandsMutex.Lock()
	s.InitTerminal()
	s.commandList = CommandList.NewCommandList()
	return s
}

func (s *Session) InitTerminal() {
	fmt.Printf("ID here = %v\n", s.ID)
	s.Terminal = Terminal.InitTerminal(s.ID)
	s.Terminal.Run()
}

func (s *Session) InitCommand() {

}
