package ProcessStatus

import (
	"sssh_server/CustomUtils"
	"sssh_server/Modules/SSH"
	"sssh_server/SessionModules/API"
	"strings"
)

type ProcessStatusModule struct {
	user string
}

type ProcessStatusModel struct {
	User     string `json:USER`
	PID      string `json:PID`
	CPU      string `json:%CPU`
	MEM      string `json:%MEM`
	VSZ      string `json:VSZ`
	RSS      string `json:RSS`
	Terminal string `json:TT`
	Stat     string `json:STAT`
	Started  string `json:STARTED`
	Time     string `json:TIME`
	Command  string `json:COMMAND`
}

func (ps *ProcessStatusModule) OnNewSession(session API.TerminalSessionInterface) {
	ps.user = session.GetUsername()
}
func (ps *ProcessStatusModule) OnNewConnection(sshSession *SSH.SSHSession) {}

func (ps *ProcessStatusModule) Close() {}

func Fields(line string, headers []string) []string {
	fields := strings.Fields(line)
	validFields := fields[:len(headers)]

	remainingFields := fields[len(headers):]

	for _, field := range remainingFields {
		validFields[len(validFields)-1] += " " + field
	}
	return validFields
}

func (ps *ProcessStatusModule) getProcessStatus() *[]map[string]string {
	res := CustomUtils.ExecuteCommandOnce("ps aux", ps.user)
	lines := strings.Split(res, "\n")
	headers := strings.Fields(lines[0])

	lines = lines[1 : len(lines)-1] // Last line is only "" so we ignore it
	process := make([]map[string]string, len(lines))

	for i := 0; i < len(lines); i++ {
		(process)[i] = make(map[string]string)
	}

	for i, line := range lines {
		if line == "" {
			continue
		}
		fields := Fields(line, headers)
		for j, field := range fields {
			(process)[i][headers[j]] = field
		}
	}

	return &process
}

var _ API.Module = (*ProcessStatusModule)(nil)
