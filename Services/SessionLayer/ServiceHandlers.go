package SessionLayer

import (
	"encoding/json"
	"fmt"
	"github.com/creack/pty"
	"io"
	"io/ioutil"
	"sssh_server/CustomUtils"
	"sssh_server/Services/CommandExecuter"
	"sssh_server/Services/CommandList"
	"sssh_server/Services/GlobalVariables"
	"sssh_server/Services/RecentCommands"
)

var commandExecuter CommandExecuter.CommandExecuter

func (s *SessionService) InitHandlers() {

	// Single run commands
	// They make a request and gets one answers
	s.HandleFunc("commandlist", s.getCommandListCommand)
	s.HandleFunc("exec", s.execCommand)
	s.HandleFunc("man", s.manCommand)
	s.HandleFunc("echo", s.echoCommand)
	s.HandleFunc("globalVars", s.globalVars)
	s.HandleFunc("setVar", s.setVar)

	// Connections
	// open connections that received and send data many times
	s.HandleFunc("commands", s.commandsConnection)
	s.HandleFunc("terminal", s.terminalConnection)
	s.HandleFunc("terminal.resize", s.terminalResizeConnection)
	s.HandleFunc("echoConnection", s.echoConnection)
}

func (service *SessionService) echoCommand(s *TerminalSession, w io.Writer, r io.Reader) {
	b, _ := CustomUtils.Read(r)
	fmt.Println("echoing : " + string(b))
	_, e := w.Write(b)
	CustomUtils.CheckPrint(e)
}

// Executes a command using the ssh connection as god intended to be
func (service *SessionService) execCommand(s *TerminalSession, w io.Writer, r io.Reader) {
	b, err := CustomUtils.Read(r)
	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	_, _ = w.Write([]byte(commandExecuter.ExecuteCommand(string(b))))
}

// Return the list of existing commands
func (service *SessionService) getCommandListCommand(s *TerminalSession, w io.Writer, r io.Reader) {
	commands := CommandList.NewCommandList()
	str, err := json.Marshal(commands.GetCommandsList())
	CustomUtils.CheckPrint(err)
	_, _ = w.Write([]byte(str))
}

// Return the manual of a given command
func (service *SessionService) manCommand(s *TerminalSession, w io.Writer, r io.Reader) {
	b, err := CustomUtils.Read(r)
	CustomUtils.CheckPrint(err)
	//fmt.Println("exec command: " + string(b))
	manCmnd := fmt.Sprintf("man %s | col -b", string(b))
	_, _ = w.Write([]byte(commandExecuter.ExecuteCommand(manCmnd)))
}

// Return an array of vars with the global variables
func (service *SessionService) globalVars(s *TerminalSession, w io.Writer, r io.Reader) {
	_, _ = w.Write([]byte(s.GlobalVars.GetVariables()))
}

// Set a given var
func (service *SessionService) setVar(s *TerminalSession, w io.Writer, r io.Reader) {
	b, _ := CustomUtils.Read(r)
	var bashVar GlobalVariables.BashVar
	_ = json.Unmarshal(b, &bashVar)
	_ = s.GlobalVars.StoreVariable(bashVar)
}

// Connections

// Echo connection
func (service *SessionService) echoConnection(s *TerminalSession, w io.Writer, r io.Reader) {
	fmt.Println("Hello")
	_, _ = io.Copy(w, r)
}

// On new commands connection opened
// continuously send the new commands
func (service *SessionService) commandsConnection(s *TerminalSession, w io.Writer, r io.Reader) {
	userSession := s
	userSession.recentCommands = RecentCommands.NewRecentCommands(userSession.Terminal, func(commands []RecentCommands.Command) {
		// TODO: properly handle error
		newCommandsJson, _ := json.Marshal(commands)
		_, _ = w.Write(newCommandsJson)
	})
	//userSession.recentCommandsMutex.Unlock()
}

// On terminal connection opened
// Simple redirect the channels
func (service *SessionService) terminalConnection(s *TerminalSession, w io.Writer, r io.Reader) {
	go func() { _, _ = io.Copy(s.Terminal, r) }()
	_, _ = io.Copy(w, s.Terminal)
}

// On terminal.resize connection opened
func (service *SessionService) terminalResizeConnection(s *TerminalSession, w io.Writer, r io.Reader) {
	var resize pty.Winsize
	b, err := ioutil.ReadAll(r)
	err = json.Unmarshal(b, &resize)
	if err != nil {
		s.Terminal.SetSize(&resize)
	}
}
