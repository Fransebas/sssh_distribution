package Sockets

import (
	"github.com/gorilla/websocket"
)

type SocketReadWriter struct {
	Socket *websocket.Conn
}

func (srw SocketReadWriter) Write(p []byte) (n int, err error) {
	err = srw.Socket.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

func (srw SocketReadWriter) Read(p []byte) (n int, err error) {
	_, p, err = srw.Socket.ReadMessage()
	return
}
