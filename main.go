package main

import (
	"flag"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/RPC"
	"sssh_server/Services/SSH"
	"sssh_server/Services/SessionLayer"
	"strings"
)

var modePtr = flag.String("mode", "server", `Select a mode for the program, available modes are: 
	server : running the sssh server
	prompt : system only function (the user shouldn't use it), it send a request to the server indicating the user typed a command, it should be use it conjunction with userid `)
var userIdPtr = flag.String("userid", "error", "Send the id of the user should be used with the mode flag set to prompt")
var historyPtr = flag.String("history", "error", "The history of the bash, should be used with the model flag set to prompt")

var portPtr = flag.Int("port", 2000, "Port for the http server")
var rpcPortPtr = flag.Int("rpcport", 2001, "Select a port for the rpc (internal process communication)")

var sessionService *SessionLayer.SessionService

var rpc *RPC.RPC

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
	fmt.Println("SSSH_USER = " + id)
	sessionService.AddCommand(string(b), id)
	//rpc.OnCommand(string(b))
}

func variables(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	CustomUtils.CheckPrint(err)

	id := r.URL.Query().Get("SSSH_USER")
	fmt.Println("SSSH_USER = " + id)
	sessionService.UpdateVariables(string(b), id)
	//rpc.OnCommand(string(b))
}

//func rpcServer() {
//	go rpc.Serve()
//}

func server() {
	SSH.GenerateNewECSDAKey()
	//r := mux.NewRouter()
	mux := http.NewServeMux()
	sessionService = SessionLayer.Constructor()
	// needed http
	mux.HandleFunc("/newcommand", newCommand)
	mux.HandleFunc("/variables", newCommand)

	log.Printf("Serving at localhost:%v...\n", (*portPtr))
	handler := cors.Default().Handler(mux)
	go sessionService.Serve()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", (*portPtr)), handler))
}

// I don't like this solution but time will tell
func updateVariables() {
	exec := CommandExecuter.CommandExecuter{}
	data := exec.ExecuteCommand("env")
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/newcommand?SSSH_USER=%v", *userIdPtr), "text/html", strings.NewReader(data))
	CustomUtils.CheckPrint(e)
}

func prompt() {
	_, e := http.Post(fmt.Sprintf("http://localhost:2000/variables?SSSH_USER=%v", *userIdPtr), "text/html", strings.NewReader(*historyPtr))
	CustomUtils.CheckPrint(e)
	updateVariables()
}

func main() {
	flag.Parse()
	rpc = RPC.New(*rpcPortPtr)
	for _, service := range SessionLayer.CommandServices {
		rpc.AddService(service)
	}

	if *modePtr == "server" {
		//rpcServer()
		server()
	} else if *modePtr == "prompt" {
		prompt()
	}
}
