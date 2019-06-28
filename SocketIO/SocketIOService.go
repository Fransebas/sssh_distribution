package SocketIO

import (
	"fmt"
	"github.com/googollee/go-engine.io"
	"github.com/googollee/go-socket.io"
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
	SocketConn *socketio.Conn
	Terminal   Terminal.Terminal
}

type SocketIOService struct {
	Server   *socketio.Server
	Sessions map[string]*Session
}

//type SockerIOParams struct {
//	port string
//}

func Constructor() (s SocketIOService) {

	opt := engineio.Options{
		//PingInterval : time.Millisecond*100,
		//PingTimeout : time.Minute,
	}

	server, e := socketio.NewServer(&opt)
	s.Server = server
	s.Sessions = make(map[string]*Session)
	CustomUtils.CheckPanic(e, "Could not initiate socket io server")

	server.OnConnect("/", func(socket socketio.Conn) error {
		socket.SetContext("")
		fmt.Println("connected:", socket.ID())
		if _, ok := s.Sessions[socket.ID()]; !ok {
			session := newSession(&socket)
			s.Sessions[socket.ID()] = session
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

func newSession(socket *socketio.Conn) (s *Session) {
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
	s.Server.OnEvent("/", "terminal", func(conn socketio.Conn, msg string) {
		fmt.Println("data: " + msg)
		s.Sessions[conn.ID()].Terminal.Write([]byte(msg))
	})

	s.Server.OnEvent("/", "command", func(conn socketio.Conn, msg string) {
		// session := s.Sessions[conn.ID()]
	})
}
