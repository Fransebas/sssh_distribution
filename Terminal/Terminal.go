package Terminal

import (
	"fmt"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"rgui/CustomUtils"
)

const (
	relativePath = "Assets/bashrc.sh"
)

/*
Holds all the important things that make up a terminal
Safe for copy, all types are pointers
*/
type Terminal struct {
	ptmx  *os.File
	ch    *chan os.Signal
	state *terminal.State
	//srw *io.ReadWriter
}

// Constructor
func InitTerminal() *Terminal {
	var t Terminal
	t.ptmx = initInteractive()
	return &t
}

// Write data into the terminal/bash
func (t *Terminal) Write(b []byte) {
	t.ptmx.Write(b)
}

/*
Set a writer that would continuous read from the terminal

We can't simple read because the read is async so we need to read all the time, hence the io.Copy
the socket or whatever interface that comunicates with it need to implement the io.Writer Interface

Don't know if this should use the "go func" or the calling function 🤷🏼‍♂️
*/
func (t *Terminal) ContinuousRead(writer io.Writer) {
	go func() { _, _ = io.Copy(writer, t.ptmx) }()
}

/*
For now this function initialize the terminal and checks continuosly if the screen change size but this doesn't work
already because the socket doesn't send that info yet
*/
func (t *Terminal) Run() {
	var err error
	t.state, err = terminal.MakeRaw(0)
	CustomUtils.CheckPanic(err, "Couldn't make terminal")

	// initialCmds(ptmx)

	//

	// Handle pty size.
	//ch := make(chan os.Signal, 1)
	//t.ch = &ch
	//signal.Notify(*t.ch, syscall.SIGWINCH)

	var winSize pty.Winsize
	winSize.X = 24
	winSize.Y = 24
	winSize.Cols = 480
	winSize.Rows = 336

	if err := pty.Setsize(t.ptmx, &winSize); err != nil {
		log.Printf("error resizing pty: %s", err)
	}
	//go func() {
	//	for range *t.ch {
	//		if err := pty.InheritSize(os.Stdin, t.ptmx); err != nil {
	//			log.Printf("error resizing pty: %s", err)
	//		}
	//	}
	//}()
	//*t.ch <- syscall.SIGWINCH // Initial resize.

	//reader := (*t.srw).(io.Reader)
	//writer := (*t.srw).(io.Writer)

	// Copy stdin to the pty and the pty to stdout.
	//go func() {_, _ = io.Copy(t.ptmx, reader) }()
	//_, _ = io.Copy(writer, t.ptmx)
}

func (t *Terminal) Close() {
	// Make sure to close the pty at the end.
	// Best effort.// Set stdin in raw mode.
	defer func() { _ = t.ptmx.Close() }()

	defer func() { _ = terminal.Restore(0, t.state) }() // Best effort.
}

func initInteractive() *os.File {
	// Handle pty size.
	relativePath := "Assets/bashrc"
	path, err := filepath.Abs(relativePath)
	fmt.Println("path " + path)
	c := exec.Command("bash", "--rcfile", path, "-i")
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
