package main

import (
	"flag"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Services/SSH"
	"sssh_server/Services/SessionLayer"
)

var port = flag.Int("port", 2000, "Select a port")

var sessionService *SessionLayer.SessionService

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
	sessionService.AddCommand(string(b), id)
}

func main() {
	flag.Parse()
	// var s ssh.SSHSession

	SSH.GenerateNewECSDAKey()
	//r := mux.NewRouter()
	mux := http.NewServeMux()
	sessionService = SessionLayer.Constructor()
	// needed http
	mux.HandleFunc("/newcommand", newCommand)

	log.Printf("Serving at localhost:%v...\n", (*port))
	handler := cors.Default().Handler(mux)
	go sessionService.Serve()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", (*port)), handler))
}
