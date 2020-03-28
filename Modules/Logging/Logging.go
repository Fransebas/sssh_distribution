package Logging

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// I know we shouldn't use caps for this but can't help it
const (
	INFO  = iota
	ERROR = iota
	DEBUG = iota
)

type Logging struct {
	infoFile  *os.File
	errorFile *os.File
}

func getBasePath() string {
	if runtime.GOOS == "windows" {
		// running no Windows
		panic("Windows not supported yet =(")
	}

	if runtime.GOOS == "darwin" {
		// running no MacOS
		return fmt.Sprintf("/var/log/sssh_server")
	} else {
		// linux and something else
		// TODO: validate other OS's
		return fmt.Sprintf("/var/log/sssh_server")
	}
}

func _exec(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	_, _ = c.Output()
}

func New() *Logging {
	l := new(Logging)

	var err error
	basePath := getBasePath()
	//os.Mkdir(basePath, 0777)

	infoName := "/info"
	errorName := "/error"

	// Fucking OS doesn't work with l.infoFile, err = os.OpenFile(basePath+infoName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)

	_exec(`mkdir`, basePath)

	_exec(`chmod`, "777", basePath)

	_exec(`touch`, basePath+infoName)
	_exec(`touch`, basePath+errorName)

	_exec("chmod", "777", basePath+infoName)
	_exec("chmod", "777", basePath+errorName)

	// If the file doesn't exist, create it, or append to the file
	l.infoFile, err = os.OpenFile(basePath+infoName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Couldn't create info log file")
		fmt.Println(err)
	}

	l.errorFile, err = os.OpenFile(basePath+errorName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Couldn't create error log file")
		fmt.Println(err)
	}

	return l
}

func (l *Logging) print(str string, level string, file *os.File) {
	t := time.Now()
	// Print every line starting with &^%$ for easy processing
	s := fmt.Sprintf("%v | %v | %v | %v \n", "&^%$", t.String(), level, str)

	//bw := bufio.NewWriter(file)

	if file == nil {
		fmt.Println("File doesn't exist")
		fmt.Println(s)
	} else {
		_, e := file.Write([]byte(s))
		if e != nil {
			fmt.Println("Error writing to file")
			fmt.Println(s)
		}
	}

}

func (l *Logging) Println(level int, str string) {
	if level == INFO {
		l.print(str, "INFO", l.infoFile)
	} else if level == DEBUG {
		l.print(str, "DEBUG", l.infoFile)
	} else {
		l.print(str, "ERROR", l.errorFile)
	}
}

func (l *Logging) Printlnf(level int, str string, vars ...interface{}) {
	s := fmt.Sprintf(str, vars)
	l.Println(level, s)
}
