package RecentCommands

import (
	"bufio"
	"io/ioutil"
	"sssh_server/CustomUtils"
	"sssh_server/Terminal"
	"strconv"
	"strings"
)

type RecentCommands struct {
	Cmnds             []Command
	terminal          *Terminal.Terminal
	OnCommandsUpdateF func([]Command)
}

type Command struct {
	Index int
	Cmnd  string
}

var lastCmd = "history"

func readFirstInt(s string) (r, charCount int) {
	ss := strings.Split(s, " ")
	r, _ = strconv.Atoi(ss[0])
	charCount = len(ss[0])
	return
}

func Constructor(terminal *Terminal.Terminal) *RecentCommands {
	rc := new(RecentCommands)
	rc.Cmnds = []Command{}
	rc.terminal = terminal
	rc.ReadCommandsInFile()
	return rc
}

func NewRecentCommands(terminal *Terminal.Terminal, f func([]Command)) *RecentCommands {
	rc := new(RecentCommands)
	rc.Cmnds = []Command{}
	rc.terminal = terminal
	rc.OnCommandsUpdateF = f
	rc.ReadCommandsInFile()
	return rc
}

func processFiles(s string) string {
	lines := strings.Split(s, "\n")
	out := ""
	i := 0
	for _, l := range lines {
		out += string(i) + " " + l + "\n"
	}
	return out
}

func (rc *RecentCommands) ReadCommandsInFile() {
	filepath := rc.terminal.HistoryFilePath
	data, err := ioutil.ReadFile(filepath)
	CustomUtils.CheckPrint(err)
	rc.UpdateRecentCommands(processFiles(string(data)))
}

func (rc *RecentCommands) OnCommandsUpdate(f func([]Command)) {
	rc.OnCommandsUpdateF = f
}

func CommandFromParts(s []string) Command {
	i, _ := strconv.Atoi(s[0])
	cmdStr := ""
	sep := ""
	for _, c := range s[1:] {
		cmdStr += sep + c
		sep = " "
	}
	return Command{
		Index: i,
		Cmnd:  cmdStr,
	}
}

func CommandFromStr(s string) Command {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	parts := []string{}
	for scanner.Scan() {
		parts = append(parts, scanner.Text())
	}
	return CommandFromParts(parts)
}

func (r *RecentCommands) UpdateRecentCommands(rawData string) []Command {
	cmnds := strings.Split(rawData, "\n")
	newCmnds := []Command{}
	for _, cmnd := range cmnds {
		if cmnd == "" {
			continue
		}
		// fmt.Println("cmd : " + cmnd + "9999\n")
		cmd := CommandFromStr(cmnd)
		newCmnds = append(newCmnds, cmd)

	}
	r.Cmnds = append(r.Cmnds, newCmnds...)
	if r.OnCommandsUpdateF != nil {
		r.OnCommandsUpdateF(newCmnds)
	}
	return newCmnds
}
