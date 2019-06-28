package SocketIO

import "github.com/graarh/golang-socketio"

type Connection interface {
	Emit(event string, msg string)
	ID() string
}

type GraarhConnectionWrapper struct {
	gosocketio.Channel
};

func (g *GraarhConnectionWrapper) Emit(event string, msg string) {
	_ = g.Channel.Emit(event, msg)
}

func (g *GraarhConnectionWrapper) ID() string {
	return g.Channel.Id()
}



type SocketServer interface {
	On(event string, func ())
	OnConnect()
	OnDisconnect()
}

type GraarhSocketServer struct {
	gosocketio.Server
}

func On(event string, )  {

}