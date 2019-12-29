package History

import (
	"bufio"
	"io/ioutil"
	"sssh_server/CustomUtils"
	"sssh_server/Services/API"
	"sssh_server/Services/SSH"
	"strconv"
	"strings"
)

type History struct {
	Cmnds             []Command
	OnCommandsUpdateF func([]Command)
}

type Command struct {
	Index int
	Cmnd  string
}

func (h *History) OnNewCommand(cmnd string) {
	h.updateRecentCommands(cmnd)
}

func (h *History) OnNewSession(session API.TerminalSessionInterface) {
	h.Cmnds = []Command{}
	config := session.GetConfig()
	h.readCommandsInFile(config.GetHistoryFilePath())
}

func (h *History) OnNewConnection(sshSession *SSH.SSHSession) {

}

func processFiles(s string) string {
	lines := strings.Split(s, "\n")
	out := ""
	i := 0
	for _, l := range lines {
		out += strconv.Itoa(i) + " " + l + "\n"
	}
	return out
}

func (h *History) readCommandsInFile(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	CustomUtils.CheckPrint(err)
	h.updateRecentCommands(processFiles(string(data)))
}

func commandFromParts(s []string) Command {
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

func commandFromStr(s string) Command {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	parts := []string{}
	for scanner.Scan() {
		parts = append(parts, scanner.Text())
	}
	return commandFromParts(parts)
}

func (h *History) updateRecentCommands(rawData string) []Command {
	cmnds := strings.Split(rawData, "\n")
	newCmnds := []Command{}
	for _, cmnd := range cmnds {
		if cmnd == "" {
			continue
		}
		// fmt.Println("cmd : " + cmnd + "9999\n")
		cmd := commandFromStr(cmnd)
		newCmnds = append(newCmnds, cmd)

	}
	h.Cmnds = append(h.Cmnds, newCmnds...)
	if h.OnCommandsUpdateF != nil {
		h.OnCommandsUpdateF(newCmnds)
	}
	return newCmnds
}

//var _ API.Service = History{}       // Verify that T implements I.
var _ API.Service = (*History)(nil)        // Verify that *T implements I.
var _ API.HistoryService = (*History)(nil) // Verify that *T implements I.
