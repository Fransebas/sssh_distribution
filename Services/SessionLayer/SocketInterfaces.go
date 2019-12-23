package SessionLayer

import (
	"github.com/desertbit/glue"
	"github.com/graarh/golang-socketio"
	"net/http"
	"sssh_server/Services/SSH"
)

type Connection interface {
	Emit(event string, msg string)
	ID() string
	GetUser() *SSH.User
}

type GraarhConnectionWrapper struct {
	*gosocketio.Channel
}

func (g GraarhConnectionWrapper) Emit(event string, msg string) {
	_ = g.Channel.Emit(event, msg)
}

func (g GraarhConnectionWrapper) ID() string {
	return g.Channel.Id()
}
func (g GraarhConnectionWrapper) GetUser() *SSH.User {
	return nil
}

type SocketServer interface {
	On(event string, f func(c Connection, msg string))
	OnConnect(f func(c Connection, msg string))
	OnDisconnect(f func(c Connection, msg string))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Serve()
}

func NewGraarhSocketServer(g *gosocketio.Server) (s GraarhSocketServer) {
	s.Server = g
	return
}

type GraarhSocketServer struct {
	*gosocketio.Server
}

func (g GraarhSocketServer) On(event string, f func(c Connection, msg string)) {
	_ = g.Server.On(event, func(c *gosocketio.Channel, msg string) string {
		var gc GraarhConnectionWrapper
		gc.Channel = c
		f(gc, msg)
		return "OK"
	})
}

func (g GraarhSocketServer) OnConnect(f func(c Connection, msg string)) {
	g.On(gosocketio.OnConnection, f)
}

func (g GraarhSocketServer) OnDisconnect(f func(c Connection, msg string)) {
	g.On(gosocketio.OnDisconnection, f)
}

func (g GraarhSocketServer) Serve() {
	//g.Server.
}

// Glue

type GlueConnectionWrapper struct {
	*glue.Socket
}

func NewGlueConnectionWrapper(socket *glue.Socket) (s GlueConnectionWrapper) {
	s.Socket = socket
	return
}

func (g GlueConnectionWrapper) Emit(event string, msg string) {
	g.Channel(event).Write(msg)
}

func (g GlueConnectionWrapper) ID() string {
	return g.Socket.ID()
}

func (g GlueConnectionWrapper) GetUser() *SSH.User {
	return nil
}

// Glue Server

type GlueSocketServer struct {
	*glue.Server
}

func NewGlueSocketServer(g *glue.Server) (s GlueSocketServer) {
	s.Server = g
	return
}

func (g GlueSocketServer) OnConnect(f func(c Connection, msg string)) {
	g.Server.OnNewSocket(func(s *glue.Socket) {
		connection := NewGlueConnectionWrapper(s)
		f(connection, "")
	})
}

func (g GlueSocketServer) Serve() {
	//g.Server.
}

// Custom

type CustomConnectionWrapper struct {
	*CustomSocket
}

func NewCustomConnectionWrapper(socket *CustomSocket) (s CustomConnectionWrapper) {
	s.CustomSocket = socket
	return
}

func (g CustomConnectionWrapper) Emit(event string, msg string) {
	_ = g.CustomSocket.Emit(event, msg)
}

func (g CustomConnectionWrapper) ID() string {
	return g.CustomSocket.ID()
}

func (g CustomConnectionWrapper) GetUser() *SSH.User {
	return g.user
}

// Custom Server
