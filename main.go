package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"rgui/CustomUtils"
	"rgui/Terminal"
)

const (
	port = ":2000"
)

type CorseMiddleware struct {
}

func (CorseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with, Content-Type, origin, authorization, accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT, PATCH")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// return "OKOK"
	server.ServeHTTP(w, r)

}

var server *socketio.Server
var middleWare CorseMiddleware
var upgrader = websocket.Upgrader{}

func setUpSocket(so *socketio.Socket) {
	err := (*so).On("data", func(msg string) {
		fmt.Println(msg)
		_, err := os.Stdout.Write([]byte(msg))
		//_, err := os.Stdin.Write()
		CustomUtils.CheckPrint(err)
	})
	CustomUtils.CheckPrint(err)
	// srw := Sockets.CreateSocketReadWriter(so, "data", "data")
	// err = Terminal.InitTerminal(srw)
	CustomUtils.CheckPanic(err, ": panic initializing the terminal")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// conn.Re
	rf, wf, err := os.Pipe()
	err = Terminal.InitTerminal(rf, wf)
	CustomUtils.CheckPanic(err, "unable to open terminal")

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		_, _ = wf.Write(p)
		bOuput := []byte{}
		_, _ = rf.Read(bOuput)
		if err := conn.WriteMessage(websocket.TextMessage, bOuput); err != nil {
			log.Println(err)
			return
		}
	}

}

func main() {
	//var err error
	//server, err = socketio.NewServer(nil)

	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_ = server.On("connection", func(so socketio.Socket) {
	//	fmt.Println("connection")
	//	setUpSocket(&so)
	//})
	//_ = server.On("error", func(so socketio.Socket, err error) {
	//	log.Println("error:", err)
	//})

	http.HandleFunc("/socket", serveWs)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(port, nil))
}
