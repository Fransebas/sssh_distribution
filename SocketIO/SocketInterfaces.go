package SocketIO

import (
	"github.com/graarh/golang-socketio"
	"net/http"
)

type Connection interface {
	Emit(event string, msg string)
	ID() string
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
