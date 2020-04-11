package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
	"os"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"sssh_server/Modules/Logging"
	"sssh_server/Programs"
	"sssh_server/SessionModules/RPC"
	"sssh_server/SessionModules/SessionLayer"
)

var sessionService *SessionLayer.SessionService

var rpc *RPC.RPC

var config Configuration.Configuration

func init() {
	//recentCommandsSrvc.Socket = Sockets.NewSocketReadWriter()
}

// This function adds a new command to the recent commands
// This is called from http from the localhost
// The bash command detects a new command and reports it back to the server
func newCommand(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	CustomUtils.CheckPrint(err)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	id := r.URL.Query().Get("SSSH_USER")
	sessionService.AddCommand(string(b), id)
	//rpc.OnCommand(string(b))
}

func variables(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	CustomUtils.CheckPrint(err)

	_, _ = fmt.Fprintln(w, "ok")
	flusher, _ := w.(http.Flusher)
	flusher.Flush()

	id := r.URL.Query().Get("SSSH_USER")
	sessionService.UpdateVariables(string(b), id)
	//rpc.OnCommand(string(b))
}

func pwd(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	CustomUtils.CheckPrint(err)

	id := r.URL.Query().Get("SSSH_USER")
	sessionService.UpdatePWD(string(b), id)
	//rpc.OnCommand(string(b))
}

func getPublickKey(w http.ResponseWriter, r *http.Request) {
	pubKey, e := sessionService.GetPubKey()
	CustomUtils.CheckPanic(e, "Couldn't read pub key:")
	pubKeyJson, e := json.Marshal(pubKey)
	CustomUtils.CheckPanic(e, "Couldn't marshal pub key:")
	_, _ = w.Write(pubKeyJson)
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func stop(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)

	if !(ip == "localhost" || ip == "127.0.0.1") {
		_ = httpServer.Shutdown(context.Background())
	}

	sessionService.Close()
	os.Exit(0)
}

var httpServer http.Server

func server(config Configuration.Configuration) {
	//r := mux.NewRouter()
	CustomUtils.Logger.Printlnf(Logging.INFO, "Serving at localhost:%v", config.Port)

	mux := http.NewServeMux()
	sessionService = SessionLayer.Constructor(config.KeyFile, config.Port)
	// needed http
	mux.HandleFunc("/newcommand", newCommand)
	mux.HandleFunc("/variables", variables)
	mux.HandleFunc("/pwd", pwd)

	mux.HandleFunc("/stop", stop)

	mux.HandleFunc("/pubKey", getPublickKey)

	handler := cors.Default().Handler(mux)
	go sessionService.Serve()

	httpServer := http.Server{Addr: fmt.Sprintf(":%v", config.HTTPPort), Handler: handler}
	e := httpServer.ListenAndServe()

	CustomUtils.CheckPanic(e, "Server stop")
	CustomUtils.Logger.Printlnf(Logging.INFO, "Stop")
}

func main() {
	config = Configuration.Configuration{}
	config.Init()

	CustomUtils.Logger.Println(Logging.INFO, fmt.Sprint(os.Args))
	CustomUtils.Logger.Println(Logging.INFO, config.String())

	rpc = RPC.New(config.RPCPort)
	for _, service := range SessionLayer.CommandServices {
		rpc.AddService(service)
	}

	if config.Mode == "server" {
		server(config)
	} else if config.Mode == "prompt" {
		Programs.Prompt(config)
	} else if config.Mode == "keygen" {
		Programs.Keygen(config)
	} else if config.Mode == "fingerprint" {
		Programs.Fingerprint(config)
	} else if config.Mode == "stop" {
		Programs.Stop(config)
	} else {
		panic("Invalid Mode " + config.Mode)
	}
}
