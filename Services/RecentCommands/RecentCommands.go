package RecentCommands

import (
	"bufio"
	"rgui/Sockets"
	"strconv"
	"strings"
)

type RecentCommands struct {
	Cmnds  []Command
	Socket *Sockets.SocketReadWriter
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

func (r *RecentCommands) UpdateRecentCommands(rawData string) {

	cmnds := strings.Split(rawData, "\n")
	for _, cmnd := range cmnds {
		if cmnd == "" {
			continue
		}
		// fmt.Println("cmd : " + cmnd + "9999\n")

		cmd := CommandFromStr(cmnd)
		r.Cmnds = append(r.Cmnds, cmd)
	}
}
