package Sockets

import (
	"github.com/googollee/go-socket.io"
)

type SocketReadWriter struct {
	Socket       *socketio.Socket
	ReadChannel  string
	WriteChannel string
	ch           *(chan []byte)
}

func CreateSocketReadWriter(Socket *socketio.Socket, ReadChannel string, WriteChannel string) (srw SocketReadWriter) {
	srw = SocketReadWriter{
		Socket:       Socket,
		ReadChannel:  ReadChannel,
		WriteChannel: WriteChannel,
		ch:           new(chan []byte),
	}

	//_ = (*srw.Socket).On(srw.ReadChannel, func(msg string) {
	//	fmt.Println(msg)
	//	(*srw.ch) <- []byte(msg)
	//})
	return
}

func (srw SocketReadWriter) Write(p []byte) (n int, err error) {
	err = (*srw.Socket).Emit(srw.WriteChannel, p)
	return len(p), err
}

func (srw SocketReadWriter) Read(p []byte) (n int, err error) {
	panic("Testing")
	//p = <- (*srw.ch)
	//return len(p), err
}
