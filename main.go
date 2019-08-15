package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/SSH"
	"sssh_server/Services/SocketIO"
)

const (
	httpport = ":2000"
	sshport  = ":2000"
)

var upgrader = websocket.Upgrader{}
var commandExecuter CommandExecuter.CommandExecuter

var socketService *SocketIO.SocketIOService

func init() {
	//recentCommandsSrvc.Socket = Sockets.NewSocketReadWriter()

}

// This function adds a new command to the recent commands
// This is called from http from the localhost
// The bash command detects a new command and reports it back to the server
func newCommand(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	CustomUtils.CheckPrint(err)

	id := r.URL.Query().Get("SSSH_USER")
	socketService.AddCommand(string(b), id)
}

func main() {
	// var s ssh.Session

	SSH.GenerateNewECSDAKey()
	//r := mux.NewRouter()
	mux := http.NewServeMux()
	socketService = SocketIO.Constructor()
	// needed http
	mux.HandleFunc("/newcommand", newCommand)

	log.Println("Serving at localhost:2000...")
	handler := cors.Default().Handler(mux)
	go socketService.Serve()

	log.Fatal(http.ListenAndServe(httpport, handler))
}
