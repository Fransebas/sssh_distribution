package SocketIO

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"net/http"
	"rgui/CustomUtils"
	"rgui/Terminal"
)

/*
safe for copy
underlying types use pointers to data
*/
type Session struct {
	ID         string
	SocketConn *Connection
	Terminal   Terminal.Terminal
}

type SocketIOService struct {
	Server   *gosocketio.Server
	Sessions map[string]*Session
}

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s SocketIOService) {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	s.Server = server
	s.Sessions = make(map[string]*Session)

	server.On(gosocketio.OnConnection, func(socket *gosocketio.Channel) {
		fmt.Println("connected:", socket.Id())
		if _, ok := s.Sessions[socket.Id()]; !ok {
			session := newSession(socket)
			s.Sessions[socket.Id()] = session
		} else {
			fmt.Println("Connection already exist")
		}
		return nil
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
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

func newSession(socket *Connection) (s *Session) {
	s = new(Session)
	s.ID = (*socket).ID()
	s.SocketConn = socket
	s.InitTerminal()
	return s
}

func (s *Session) InitTerminal() {
	s.Terminal = *Terminal.InitTerminal()
	s.Terminal.Run()
	socketWriter := NewSocketReadWriter(s.SocketConn, "terminal")
	s.Terminal.ContinuousRead(socketWriter)
}

func (s *Session) InitCommand() {

}

func (s *SocketIOService) InitEvents() {
	_ = s.Server.On("terminal", func(conn *gosocketio.Channel, msg string) {
		fmt.Println("data: " + msg)
		s.Sessions[conn.ID()].Terminal.Write([]byte(msg))
	})

	_ = s.Server.On("command", func(conn *gosocketio.Channel, msg string) {
		// session := s.Sessions[conn.ID()]
	})
}
