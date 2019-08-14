package Terminal

import (
	"fmt"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"rgui/CustomUtils"
	"rgui/Services/SSH"
	"sync"
)

const (
	relativePath = "Assets/bashrc.sh"
)

var FILE_NAME = "testhistory" // Constant for now, but this should change per session or user

/*
Holds all the important things that make up a terminal
Safe for copy, all types are pointers
*/
type Terminal struct {
	ptmx            *os.File
	ch              *chan os.Signal
	state           *terminal.State
	resizeMux       *sync.Mutex
	user            *SSH.User
	HistoryFilePath string
	//srw *io.ReadWriter
}

// Constructor
func InitTerminal(id string) *Terminal {
	var t Terminal

	t.resizeMux = new(sync.Mutex)

	// The file base path is the Assets, maybe there is a better place like /etc or something
	basePath, err := filepath.Abs("Assets/")
	CustomUtils.CheckPanic(err, "Could not create history file for the session")
	t.HistoryFilePath = fmt.Sprintf("%v/%v", basePath, FILE_NAME)

	t.ptmx = initInteractive(id, t.HistoryFilePath)
	return &t
}

// Write data into the terminal/bash
func (t *Terminal) Write(b []byte) (int, error) {
	return t.ptmx.Write(b)
}

func (t *Terminal) Read(b []byte) (int, error) {
	return t.ptmx.Read(b)
}

/*
Set a writer that would continuous read from the terminal

We can't simple read because the read is async so we need to read all the time, hence the io.Copy
the socket or whatever interface that comunicates with it need to implement the io.Writer Interface

Don't know if this should use the "go func" or the calling function 🤷🏼‍♂️
*/
//func (t *Terminal) ContinuousRead(writer io.Writer) {
//	go func() { _, _ = io.Copy(writer, t.ptmx) }()
//}

/*
For now this function initialize the terminal and checks continuosly if the screen change size but this doesn't work
already because the socket doesn't send that info yet
*/
func (t *Terminal) Run() {
	var err error
	t.state, err = terminal.MakeRaw(0)
	CustomUtils.CheckPanic(err, "Couldn't make terminal")

	var winSize pty.Winsize
	winSize.X = 24
	winSize.Y = 24
	winSize.Cols = 480
	winSize.Rows = 336

	if err := pty.Setsize(t.ptmx, &winSize); err != nil {
		log.Printf("error resizing pty: %s", err)
	}
}

func (t *Terminal) SetSizeVals(X, Y, COLS, ROWS uint16) {
	var winSize pty.Winsize
	winSize.X = X
	winSize.Y = Y
	winSize.Cols = COLS
	winSize.Rows = ROWS
	t.SetSize(&winSize)
}

func (t *Terminal) SetSize(winSize *pty.Winsize) {
	t.resizeMux.Lock()
	defer t.resizeMux.Unlock()
	if err := pty.Setsize(t.ptmx, winSize); err != nil {
		log.Printf("error resizing pty: %s", err)
	}
}

func (t *Terminal) Close() {
	// Make sure to close the pty at the end.
	// Best effort.// Set stdin in raw mode.
	defer func() { _ = t.ptmx.Close() }()
	defer func() { _ = terminal.Restore(0, t.state) }() // Best effort.
}

func initInteractive(ID string, historyPath string) *os.File {
	// Handle pty size.
	basePath, err := filepath.Abs("Assets/")
	fmt.Printf("basePath = %v \n", basePath)
	relativePath := fmt.Sprintf("%v/bashrc", basePath)
	path := relativePath
	fmt.Println("path " + path)

	// Send the initialization file the variables it's going to use
	bash := fmt.Sprintf("export SSSH_USER=%v; export HIST_FILE_NAME=%v; bash --rcfile %s -i ", ID, historyPath, path)
	c := exec.Command("bash", "-c", bash)
	//c := exec.Command("bash", "--rcfile", path, "-i")
	CustomUtils.CheckPanic(err, "Could not initialize Terminal")
	// Start the command with a pty.

	ptmx, err := pty.Start(c)
	CustomUtils.CheckPanic(err, "Could not initialize Terminal")

	return ptmx
}

//
//func initLogin() *os.File{
//	// Handle pty size.
//	c := exec.Command("bash", "-l")
//	// Start the command with a pty.
//	ptmx, err := pty.Start(c)
//	initialCmds(ptmx)
//	CustomUtils.CheckPanic(err, "Could not initialize Terminal");
//	return ptmx
//}
//
//func ReadNPrint(pty *os.File){
//	r := []byte{}
//	_, _ = pty.Read(r)
//	fmt.Println("res :" +string(r))
//}
//
//func initialCmds(pty *os.File) {
//	path, err := filepath.Abs(relativePath)
//	CustomUtils.CheckPrint(err)
//	fmt.Println("path "+ path)
//	_, _ = pty.WriteString(path + "\n")
//}

var ignoreInternalCmds = " export HISTCONTROL=ignorespace ; history -d $(history 1) \n\n"
var readCmdHistory = func(cmd string) string { return fmt.Sprintf(" export PROMPT_COMMAND='%s' \n\n", cmd) }
var runOnCommand = func(cmd string) string { return fmt.Sprintf(" export PROMPT_COMMAND='%s' \n\n", cmd) }
var getLastCommand = func(cmd string) string { return fmt.Sprintf(" history 1 | %s", cmd) }
var makeRqst = func(n string) string {
	return fmt.Sprintf("curl -d \"$(history %s)\" http://localhost:2000/newcommand", n)
}
