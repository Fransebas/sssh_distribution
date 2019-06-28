package SocketIO

import (
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"net/http"
	"rgui/Terminal"
)

/*
safe for copy
underlying types use pointers to data
*/
type Session struct {
	ID         string
	SocketConn Connection
	Terminal   Terminal.Terminal
}

type SocketIOService struct {
	Server   SocketServer
	Sessions map[string]*Session
}

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s SocketIOService) {

	server := NewGraarhSocketServer(gosocketio.NewServer(transport.GetDefaultWebsocketTransport()))
	s.Server = server
	s.Sessions = make(map[string]*Session)

	server.On(gosocketio.OnConnection, func(socket Connection, msg string) {
		fmt.Println("connected:", socket)
		if _, ok := s.Sessions[socket.ID()]; !ok {
			session := newSession(socket)
			s.Sessions[socket.ID()] = session
		} else {
			fmt.Println("Connection already exist")
		}
	})
	server.OnDisconnect(func(s Connection, msg string) {
		fmt.Println("closed", msg)
	})
	s.InitEvents()

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

func newSession(socket Connection) (s *Session) {
	s = new(Session)
	s.ID = socket.ID()
	s.SocketConn = socket
	s.InitTerminal()
	return s
}

func (s *Session) InitTerminal() {
	s.Terminal = *Terminal.InitTerminal()
	s.Terminal.Run()
	socketWriter := NewSGnrSocketIOWriter(s.SocketConn, "terminal")
	s.Terminal.ContinuousRead(socketWriter)
}

func (s *Session) InitCommand() {

}

func (s *SocketIOService) InitEvents() {
	s.Server.On("terminal", func(socket Connection, msg string) {
		fmt.Println("data: " + msg)
		s.Sessions[socket.ID()].Terminal.Write([]byte(msg))
	})

	s.Server.On("command", func(socket Connection, msg string) {
		// session := s.Sessions[conn.ID()]
		fmt.Println("command data: " + msg)
	})
}
