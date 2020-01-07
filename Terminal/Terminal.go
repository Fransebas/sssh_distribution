package Terminal

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sssh_server/CustomUtils"
	"sssh_server/Services/SSH"
	"sync"
)

const (
	relativePath = "Assets/bashrc.sh"
)

/*
Holds all the important things that make up a terminal
Safe for copy, all types are pointers
*/
type Terminal struct {
	ptmx      *os.File
	ch        *chan os.Signal
	state     *terminal.State
	resizeMux *sync.Mutex
	user      *SSH.User
	tmp       *os.File
	buffer    *CustomUtils.FixedDeque
	//srw *io.ReadWriter
}

func (t *Terminal) GetReader() TerminalReader {
	tr := TerminalReader{}
	tr.buffer = t.buffer
	tr.offset = 0
	tr.terminal = t
	return tr
}

func (t *Terminal) GetBuffer() []byte {
	return t.buffer.Bytes()
}

// Constructor
func InitTerminal(id string, historyPath string) *Terminal {
	var t Terminal

	t.buffer = CustomUtils.New(1000000) /// 1000000 is 1 MB maybe I should use less
	t.resizeMux = new(sync.Mutex)

	// debug
	var err error
	t.tmp, err = os.Create("tmp2")
	if err != nil {
		panic(err)
	}
	// end

	// The file base path is the Assets, maybe there is a better place like /etc or something
	//basePath, err := filepath.Abs("Assets/")
	//CustomUtils.CheckPanic(err, "Could not create history file for the session")
	//t.HistoryFilePath = fmt.Sprintf("%v/%v", basePath, FILE_NAME)

	t.ptmx = initInteractive(id, historyPath)
	return &t
}

// Write data into the terminal/bash
func (t *Terminal) Write(b []byte) (int, error) {
	return t.ptmx.Write(b)
}

func (t *Terminal) Read(b []byte) (int, error) {
	n, e := t.ptmx.Read(b)
	for i := 0; i < n; i++ {
		t.buffer.Insert(b[i])
	}
	return n, e
}

func (t *Terminal) read() {
	b := make([]byte, 32*1024)
	n, _ := t.ptmx.Read(b)
	for i := 0; i < n; i++ {
		t.buffer.Insert(b[i])
	}
	t.tmp.Write(b)
}

/*
Set a writer that would continuous read from the terminal

We can't simple read because the read is async so we need to read all the time, hence the io.Copy
the socket or whatever interface that comunicates with it need to implement the io.Writer Interface

Don't know if this should use the "go func" or the calling function 🤷🏼‍♂️
*/
//func (t *TerminalService) ContinuousRead(writer io.Writer) {
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

// Creates and interactive bash session based on the ID and the file to use as history file, here we add the file to run for the initialization,
// i.e. we change the bashrc for our own version that by itself call the user bashrc
// the init file can be found in `/sssh_server/Assets/bashrc` *this should change to a local directory*
func initInteractive(ID string, historyPath string) *os.File {
	// Handle pty size.
	basePath, err := filepath.Abs("Assets/")
	fmt.Printf("basePath = %v \n", basePath)
	relativePath := fmt.Sprintf("%v/bashrc", basePath)
	path := relativePath
	fmt.Println("path " + path)

	// Send the initialization file the variables it's going to use
	// TODO: fix the absolute path thingy
	_ = os.Setenv("SYMBIONT", "~/go/src/symbiont/main.go")
	bash := fmt.Sprintf("export SSSH=%v; export SSSH_USER=%v; export HIST_FILE_NAME=%v; bash --rcfile %s -i ", "~/go/src/sssh_server/sssh_server", ID, historyPath, path)
	c := exec.Command("bash", "-c", bash)

	//c := exec.Command("bash", "--rcfile", path, "-i")
	//CustomUtils.CheckPanic(err, "Could not initialize TerminalService")
	// Start the command with a pty.

	ptmx, err := pty.Start(c)
	CustomUtils.CheckPanic(err, "Could not initialize TerminalService")

	return ptmx
}

//
//func initLogin() *os.File{
//	// Handle pty size.
//	c := exec.Command("bash", "-l")
//	// Start the command with a pty.
//	ptmx, err := pty.Start(c)
//	initialCmds(ptmx)
//	CustomUtils.CheckPanic(err, "Could not initialize TerminalService");
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
