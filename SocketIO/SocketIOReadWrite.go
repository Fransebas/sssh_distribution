package SocketIO

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"sync"
)

/*

This interface is because the terminal needs to read and write constantly to the socket
and it needs to implement the Read and Write interface
*/
type SocketIOWriter struct {
	Socket            *socketio.Conn
	SocketBufferMutex *sync.Mutex
	event             string
}

/*
Constructor for the SocketIOReadWriter,
the buffers are for more complex behavior but I think it wouldn't be needed
*/

func NewSocketReadWriter(Socket *socketio.Conn, event string) *SocketIOWriter {
	sw := SocketIOWriter{
		SocketBufferMutex: new(sync.Mutex),
		Socket:            Socket,
		event:             event,
	}
	return &sw
}

//func (srw *SocketIOReadWriter) AsyncWrite(p []byte) {
//	fmt.Println("message 2 : " + string(p))
//	*srw.SocketBft <- p
//}

func (srw *SocketIOWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Output = %s \n", string(p))
	//srw.SocketBufferMutex.Lock()
	//defer srw.SocketBufferMutex.Unlock()
	//defer fmt.Println("UNLOCKED")
	(*srw.Socket).Emit(srw.event, string(p))
	(*srw.Socket).Emit("test", "")
	return len(p), err
}

//func (srw *SocketIOReadWriter) Read(p []byte) (n int, err error) {
//	_, b, err := srw.Socket.ReadMessage()
//	for i, d := range b{
//		p[i] = d
//	}
//	return len(b), err
//}

type GnrSocketIOWriter struct {
	Socket            Connection
	SocketBufferMutex *sync.Mutex
	event             string
}

func NewSGnrSocketIOWriter(Socket Connection, event string) *GnrSocketIOWriter {
	sw := GnrSocketIOWriter{
		SocketBufferMutex: new(sync.Mutex),
		Socket:            Socket,
		event:             event,
	}
	return &sw
}

func (srw *GnrSocketIOWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Output = %s \n", string(p))
	//srw.SocketBufferMutex.Lock()
	//defer srw.SocketBufferMutex.Unlock()
	//defer fmt.Println("UNLOCKED")
	srw.Socket.Emit(srw.event, string(p))
	return len(p), err
}
