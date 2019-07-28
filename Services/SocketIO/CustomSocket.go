package SocketIO

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"rgui/CustomUtils"
	"rgui/Services/SSH"
	"sync"
	"time"
)

//type SocketServer interface {
//	On(event string, f func(c Connection, msg string))
//	OnConnect(f func(c Connection, msg string))
//	OnDisconnect(f func(c Connection, msg string))
//	ServeHTTP(w http.ResponseWriter, r *http.Request)
//	Serve()
//}
//
//type Connection interface {
//	Emit(event string, msg string)
//	ID() string
//}

type CustomSocketServer struct {
	f        func(*CustomSocket)
	sessions map[string]*CustomSocket
}

type CustomSocket struct {
	connection *websocket.Conn
	events     map[string]func(msg string)
	writeMutex *sync.Mutex
	id         string
	close      func(msg string)
	isOpen     bool
	user       *SSH.User
	isAuth     bool
	server     *CustomSocketServer
}

type Message struct {
	Event     string `json:"event"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
	ID        string `json:"ID"`
}

// This should hold all kind of security keys and stuff
type SessionInfo struct {
	ID string
}

var upgrader = websocket.Upgrader{}

func NewCustomSocketServer() (c *CustomSocketServer) {
	c = new(CustomSocketServer)
	c.sessions = make(map[string]*CustomSocket)
	return c
}

func NewCustomSocket(conn *websocket.Conn, server *CustomSocketServer) (s *CustomSocket) {
	s = new(CustomSocket)
	s.writeMutex = new(sync.Mutex)
	s.events = make(map[string]func(msg string))
	s.connection = conn
	s.isOpen = true
	s.user = SSH.NewUser()
	s.isAuth = false
	s.server = server
	return
}

func (s *CustomSocket) copySocket(socket *CustomSocket) {
	socket.close = s.close
	socket.events = s.events
}

func (c *CustomSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) *CustomSocket {
	connection, err := upgrader.Upgrade(w, r, nil)
	socket := NewCustomSocket(connection, c)
	CustomUtils.CheckPrint(err)

	socket.On("session", func(msg string) {
		fmt.Printf("new session %v \n", msg)
		// TODO: the security here is completely dummy!!!!
		var sessionInfo SessionInfo
		err := json.Unmarshal([]byte(msg), &sessionInfo)
		// TODO: figure out what to do with the error
		err = SSH.ProofClaims([]byte(sessionInfo.ID))
		if err == nil {
			socket.id = sessionInfo.ID
			socket.user.ID = sessionInfo.ID
			_ = socket.Emit("session", "ACK")
			// the user is telling the truth
			socket.isAuth = true

			if prevSocket, ok := c.sessions[socket.id]; ok {
				// copy all the functions of the prev socket to the new socket
				prevSocket.copySocket(socket)
				c.sessions[sessionInfo.ID] = socket
			} else {
				if c.f != nil {
					c.f(socket)
				}
			}
			c.sessions[sessionInfo.ID] = socket
		} else {
			// he is a filthy liar
			socket.isOpen = false
		}
	})
	go socket.run()
	return socket
}

func (c *CustomSocketServer) OnConnect(f func(*CustomSocket)) {
	c.f = f
}

func (c *CustomSocketServer) Serve() {
}

func (s *CustomSocket) run() {
	for s.isOpen {

		mst, msg, err := s.connection.ReadMessage()
		//
		//if s.isAuth {
		//	// This check is only after the socket has presented the credentials to authenticate
		//}

		if err != nil || mst == -1 || mst == websocket.CloseGoingAway || mst == websocket.CloseNormalClosure {
			_ = s.connection.Close()
			s.isOpen = false
			delete(s.server.sessions, s.user.ID)
			if s.close != nil {
				s.close(string(websocket.CloseGoingAway))
			}
			return
		}
		CustomUtils.CheckPrint(err)
		var message Message
		_ = json.Unmarshal(msg, &message)
		var ok bool
		var event func(msg string)
		if event, ok = s.events[message.Event]; ok {
			event(message.Msg)
		} else {
			fmt.Printf("Event %s is not registered \n", message.Event)
		}
	}
}

func (s *CustomSocket) Emit(event string, msg string) (e error) {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	var message Message
	message.Event = event
	message.Msg = msg
	message.ID = s.id
	message.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	messageText, err := json.Marshal(message)
	CustomUtils.CheckPrint(err)
	e = s.connection.WriteMessage(websocket.TextMessage, messageText)
	return
}

func (s *CustomSocket) ID() string {
	return s.id
}

func (s *CustomSocket) GetUser() *SSH.User {
	return s.user
}

func (s *CustomSocket) Write(event string, msg string) (e error) {
	//msgByte, err := SSH.SSHEncode([]byte(msg), s.user)
	//CustomUtils.CheckPrint(err)
	return s.Emit(event, msg)
}

func (s *CustomSocket) On(event string, f func(msg string)) {
	s.events[event] = f
}

func (s *CustomSocket) OnClose(f func(msg string)) {
	s.close = f
}
