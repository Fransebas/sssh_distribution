package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"rgui/CustomUtils"
	"rgui/Services/CommandExecuter"
	"rgui/Services/RecentCommands"
	"rgui/Services/SSH"
	"rgui/Services/SocketIO"
	"rgui/Terminal"
)

const (
	port = ":2000"
)

var upgrader = websocket.Upgrader{}
var recentCommandsSrvc RecentCommands.RecentCommands
var commandExecuter CommandExecuter.CommandExecuter

var terminal *Terminal.Terminal
var socketService *SocketIO.SocketIOService

func init() {
	//recentCommandsSrvc.Socket = Sockets.NewSocketReadWriter()

}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	// conn, err := upgrader.Upgrade(w, r, nil)
	//CustomUtils.CheckPrint(err)

	// recentCommandsSrvc.Socket.SetSocket(conn)

	// terminal = Terminal.InitTerminal(*recentCommandsSrvc.Socket, false)

	//CustomUtils.CheckPanic(err, "unable to open terminal")

	go func() { terminal.Run() }()
}

func serveCmdSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	//conn, err := upgrader.Upgrade(w, r, nil)
	//CustomUtils.CheckPrint(err)
	//recentCommandsSrvc.Socket.Socket = conn
}

func newCommand(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	b, err := ioutil.ReadAll(r.Body)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	CustomUtils.CheckPrint(err)

	id := r.URL.Query().Get("SSSH_USER")
	socketService.AddCommand(string(b), id)
}

func user(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	user := SSH.NewUser()
	userJson, err := json.Marshal(user)
	CustomUtils.CheckPrint(err)
	_, _ = fmt.Fprintln(w, string(userJson))
	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func getCommandList(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	// b, err := ioutil.ReadAll(r.Body)
	// CustomUtils.CheckPrint(err)

	id := r.URL.Query().Get("SSSH_USER")
	res := socketService.GetCommandList(id, "")
	_, _ = fmt.Fprintln(w, res)

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
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

func testCipher(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	res, _ := SSH.SSHEncode(b, nil)
	fmt.Fprintf(w, string(res))

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func testDecipher(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	res, _ := SSH.SSHDecode(b, nil)
	fmt.Fprintf(w, string(res))

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func testConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func man(w http.ResponseWriter, r *http.Request) {
	// todo: add security
	b, err := ioutil.ReadAll(r.Body)

	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	manCmnd := fmt.Sprintf("man %s | col -b", string(b))
	_, _ = fmt.Fprintln(w, commandExecuter.ExecuteCommand(manCmnd))

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

func TestPrintFile(path string) {
	fmt.Println("reading file ...")
	data, err := ioutil.ReadFile(path)
	CustomUtils.CheckPrint(err)
	fmt.Println(data)
}

func startSSHServer() {

}

func main() {
	// var s ssh.Session

	SSH.GenerateNewECSDAKey()
	//r := mux.NewRouter()
	mux := http.NewServeMux()
	socketService = SocketIO.Constructor()

	mux.HandleFunc("/cmdsocket", serveCmdSocket)
	mux.HandleFunc("/socket", serveWs)
	mux.HandleFunc("/newcommand", newCommand)
	mux.HandleFunc("/user", user)
	mux.HandleFunc("/commandlist", getCommandList)
	mux.HandleFunc("/exec", execCommand)
	mux.HandleFunc("/man", man)
	mux.HandleFunc("/", testConnection)

	mux.HandleFunc("/encrypt", testCipher)
	mux.HandleFunc("/decrypt", testDecipher)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:2000...")

	mux.HandleFunc("/ws", socketService.SocketIOFix)
	//mux.Handle("/socket.io/", socketService.Server)
	handler := cors.Default().Handler(mux)
	//TestPrintFile("cert.pem")
	//TestPrintFile("key.pem")
	log.Fatal(http.ListenAndServe(port, handler))
	//log.Fatal(http.ListenAndServeTLS(port, "mysitename.crt", "mysitename.key", handler))
}
