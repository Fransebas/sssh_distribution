package SocketIO

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pty"
	"net/http"
	"rgui/CustomUtils"
	"rgui/Services/CommandList"
	"rgui/Services/RecentCommands"
	"rgui/Terminal"
)

/*
safe for copy
underlying types use pointers to data
*/
type Session struct {
	ID             string
	SocketConn     Connection
	Terminal       *Terminal.Terminal
	recentCommands *RecentCommands.RecentCommands
	commandList    *CommandList.CommandList
}

type SocketIOService struct {
	Server   *CustomSocketServer
	Sessions map[string]*Session
}

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s *SocketIOService) {

	s = new(SocketIOService)

	server := NewCustomSocketServer()

	s.Server = server
	s.Sessions = make(map[string]*Session)

	server.OnConnect(func(socket *CustomSocket) {
		s.InitEvents(socket)
		session := newSession(NewCustomConnectionWrapper(socket))
		s.Sessions[socket.ID()] = session
	})

	go s.Server.Serve()
	return s
}

func (s *SocketIOService) SocketIOFix(w http.ResponseWriter, r *http.Request) {
	allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	}
	if r.Method == "OPTIONS" {
		return
	}
	//r.Header.Set("Origin", "");
	r.Header.Del("Origin")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	s.Server.ServeHTTP(w, r)
	//w
}

func (ss *SocketIOService) AddCommand(data string, id string) {

	fmt.Printf("data %v \n", data)
	if _, ok := ss.Sessions[id]; ok {
		newCommands := ss.Sessions[id].recentCommands.UpdateRecentCommands(data)
		fmt.Printf("Number of commands= %v \n", len(newCommands))
		//newCommandsJson, err := json.Marshal(newCommands)
		//if err == nil {
		//
		//	socket.SocketConn.Emit("commands", string(newCommandsJson))
		//}
	}
}

func (ss *SocketIOService) GetCommandList(id string, typpe string) string {
	//if socket, ok := ss.Sessions[id]; ok {
	commands := CommandList.NewCommandList()
	// commands := socket.commandList.GetCommandsList()
	str, err := json.Marshal(commands.GetCommandsList())
	CustomUtils.CheckPrint(err)
	return string(str)
	//}
	//return ""
}

func newSession(socket Connection) (s *Session) {
	fmt.Printf("New Soxket %v\n", socket.ID())
	s = new(Session)
	s.ID = socket.ID()
	s.SocketConn = socket
	s.InitTerminal()
	s.commandList = CommandList.NewCommandList()
	s.recentCommands = RecentCommands.NewRecentCommands(s.Terminal, func(commands []RecentCommands.Command) {
		newCommandsJson, err := json.Marshal(commands)
		if err == nil {
			socket.Emit("commands", string(newCommandsJson))
		}
	})
	return s
}

func (s *Session) InitTerminal() {
	fmt.Printf("ID here = %v\n", s.SocketConn.ID())
	s.Terminal = Terminal.InitTerminal(s.SocketConn.GetUser())
	s.Terminal.Run()
	socketWriter := NewSGnrSocketIOWriter(s.SocketConn, "terminal")
	s.Terminal.ContinuousRead(socketWriter)
}

func (s *Session) InitCommand() {

}

func (s *SocketIOService) InitEvents(socket *CustomSocket) {
	(*socket).On("terminal", func(data string) {
		s.Sessions[socket.ID()].Terminal.Write([]byte(data))
	})

	(*socket).On("terminal.resize", func(data string) {
		var resize pty.Winsize
		err := json.Unmarshal([]byte(data), &resize)
		if err != nil {
			s.Sessions[socket.id].Terminal.SetSize(&resize)
		}
	})
}
