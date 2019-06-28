package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"rgui/CustomUtils"
	"rgui/Services/CommandExecuter"
	"rgui/Services/RecentCommands"
	"rgui/SocketIO"
	"rgui/Terminal"
)

const (
	port = ":2000"
)

var upgrader = websocket.Upgrader{}
var recentCommandsSrvc RecentCommands.RecentCommands
var commandExecuter CommandExecuter.CommandExecuter

var terminal *Terminal.Terminal

func init() {
	//recentCommandsSrvc.Socket = Sockets.NewSocketReadWriter()

}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	CustomUtils.CheckPrint(err)

	recentCommandsSrvc.Socket.SetSocket(conn)

	// terminal = Terminal.InitTerminal(*recentCommandsSrvc.Socket, false)

	CustomUtils.CheckPanic(err, "unable to open terminal")

	go func() { terminal.Run() }()
}

func serveCmdSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	CustomUtils.CheckPrint(err)
	recentCommandsSrvc.Socket.Socket = conn
}

func newCommand(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	b, err := ioutil.ReadAll(r.Body)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	CustomUtils.CheckPrint(err)
	// fmt.Println("Command\n" + string(b))
	recentCommandsSrvc.UpdateRecentCommands(string(b))

}

func execCommand(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	b, err := ioutil.ReadAll(r.Body)

	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	_, _ = fmt.Fprintln(w, commandExecuter.ExecuteCommand(string(b)))

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func main() {
	//r := mux.NewRouter()
	mux := http.NewServeMux()

	mux.HandleFunc("/cmdsocket", serveCmdSocket)
	mux.HandleFunc("/socket", serveWs)
	mux.HandleFunc("/newcommand", newCommand)
	mux.HandleFunc("/exec", execCommand)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:2000...")
	//http.Handle("/api", r)
	socketService := SocketIO.Constructor()
	// socketService.InitEvents()
	mux.HandleFunc("/socket.io/", socketService.SocketIOFix)
	//mux.Handle("/socket.io/", socketService.Server)
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(port, handler))
}
