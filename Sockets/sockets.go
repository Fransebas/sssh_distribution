package Sockets

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type SocketReadWriter struct {
	Socket            *websocket.Conn
	SocketBft         *chan []byte
	SocketBufferMutex *sync.Mutex
	isOpen            bool
}

func NewSocketReadWriter() *SocketReadWriter {
	ch := make(chan []byte, 1000)
	sw := SocketReadWriter{
		SocketBft:         &ch,
		SocketBufferMutex: new(sync.Mutex),
		Socket:            nil,
		isOpen:            false,
	}
	go sw.async()
	return &sw
}

func (srw *SocketReadWriter) SetSocket(socket *websocket.Conn) {
	srw.Socket = socket
	srw.isOpen = true
}

func (srw *SocketReadWriter) async() {
	for b := range *srw.SocketBft {
		for !srw.isOpen {
			//fmt.Printf("Waiting for connection %v \n", srw.isOpen)
			time.Sleep(5000 * time.Millisecond)
		}
		srw.SocketBufferMutex.Lock()
		_ = srw.Socket.WriteMessage(websocket.TextMessage, b)
		srw.SocketBufferMutex.Unlock()
	}
}

func (srw *SocketReadWriter) AsyncWrite(p []byte) {
	*srw.SocketBft <- p
}

func (srw *SocketReadWriter) Write(p []byte) (n int, err error) {
	srw.SocketBufferMutex.Lock()
	defer srw.SocketBufferMutex.Unlock()
	err = srw.Socket.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

func (srw *SocketReadWriter) Read(p []byte) (n int, err error) {
	_, b, err := srw.Socket.ReadMessage()
	for i, d := range b {
		p[i] = d
	}
	return len(b), err
}
