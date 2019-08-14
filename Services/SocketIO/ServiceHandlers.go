package SocketIO

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pty"
	"io"
	"io/ioutil"
	"rgui/CustomUtils"
	"rgui/Services/CommandExecuter"
	"rgui/Services/CommandList"
	"rgui/Services/RecentCommands"
	"rgui/Services/SSH"
)

var commandExecuter CommandExecuter.CommandExecuter

func (s *SocketIOService) InitHandlers() {

	// Single run commands
	// They make a request and gets one answers
	s.Server.HandleFunc("commandlist", s.getCommandListCommand)
	s.Server.HandleFunc("exec", s.execCommand)
	s.Server.HandleFunc("man", s.manCommand)
	s.Server.HandleFunc("echo", s.echoCommand)

	// Connections
	// open connections that received and send data many times
	s.Server.HandleFunc("commands", s.commandsConnection)
	s.Server.HandleFunc("terminal", s.terminalConnection)
	s.Server.HandleFunc("terminal.resize", s.terminalResizeConnection)
	s.Server.HandleFunc("echoConnection", s.echoConnection)
}

func (service *SocketIOService) echoCommand(s *SSH.Session, w io.Writer, r io.Reader) {
	b, _ := CustomUtils.Read(r)
	fmt.Println("echoing : " + string(b))
	_, e := w.Write(b)
	CustomUtils.CheckPrint(e)
}

// Executes a command using the ssh connection as god intended to be
func (service *SocketIOService) execCommand(s *SSH.Session, w io.Writer, r io.Reader) {
	b, err := CustomUtils.Read(r)
	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	_, _ = w.Write([]byte(commandExecuter.ExecuteCommand(string(b))))
}

// Return the list of existing commands
func (service *SocketIOService) getCommandListCommand(s *SSH.Session, w io.Writer, r io.Reader) {
	commands := CommandList.NewCommandList()
	str, err := json.Marshal(commands.GetCommandsList())
	CustomUtils.CheckPrint(err)
	_, _ = w.Write([]byte(str))
}

// Return the manual of a given command
func (service *SocketIOService) manCommand(s *SSH.Session, w io.Writer, r io.Reader) {
	b, err := CustomUtils.Read(r)
	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	manCmnd := fmt.Sprintf("man %s | col -b", string(b))
	w.Write([]byte(commandExecuter.ExecuteCommand(manCmnd)))
}

// Connections

// Echo connection
func (service *SocketIOService) echoConnection(s *SSH.Session, w io.Writer, r io.Reader) {
	fmt.Println("Hello")
	_, _ = io.Copy(w, r)
}

// On new commands connection opened
// continuously send the new commands
func (service *SocketIOService) commandsConnection(s *SSH.Session, w io.Writer, r io.Reader) {
	userSession := service.Sessions[s.GetSessionID()]
	fmt.Println("3")
	userSession.recentCommands = RecentCommands.NewRecentCommands(userSession.Terminal, func(commands []RecentCommands.Command) {
		// TODO: properly handle error
		newCommandsJson, _ := json.Marshal(commands)
		_, _ = w.Write(newCommandsJson)
	})
	userSession.recentCommandsMutex.Unlock()
}

// On terminal connection opened
// Simple redirect the channels
func (service *SocketIOService) terminalConnection(s *SSH.Session, w io.Writer, r io.Reader) {
	go func() { _, _ = io.Copy(service.Sessions[s.GetSessionID()].Terminal, r) }()
	_, _ = io.Copy(w, service.Sessions[s.GetSessionID()].Terminal)
}

// On terminal.resize connection opened
func (service *SocketIOService) terminalResizeConnection(s *SSH.Session, w io.Writer, r io.Reader) {
	var resize pty.Winsize
	b, err := ioutil.ReadAll(r)
	err = json.Unmarshal(b, &resize)
	if err != nil {
		service.Sessions[s.GetSessionID()].Terminal.SetSize(&resize)
	}
}
