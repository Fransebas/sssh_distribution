package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
	"rgui/CustomUtils"
	"rgui/Sockets"
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

func setUpSocket(so *socketio.Socket) {
	err := (*so).On("data", func(msg string) {
		fmt.Println(msg)
		_, err := os.Stdout.Write([]byte(msg))
		//_, err := os.Stdin.Write()
		CustomUtils.CheckPrint(err)
	})
	CustomUtils.CheckPrint(err)
	srw := Sockets.CreateSocketReadWriter(so, "data", "data")
	err = Terminal.InitTerminal(srw)
	CustomUtils.CheckPanic(err, ": panic initializing the terminal")
}

func main() {
	var err error
	server, err = socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	_ = server.On("connection", func(so socketio.Socket) {
		fmt.Println("connection")
		setUpSocket(&so)
	})
	_ = server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", middleWare)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(port, nil))
}
