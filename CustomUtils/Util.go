package CustomUtils

import (
	"io"
	"math/rand"
	"os/exec"
	"sssh_server/Modules/Logging"
)

var Logger *Logging.Logging

func init() {
	Logger = Logging.New()
}

func CheckPanic(e error, msg string) {
	if e != nil {
		Logger.Printlnf(Logging.ERROR, "%v", e)
		panic(e.Error() + msg)
	}
}

func Print(msg string) {
	Logger.Println(Logging.INFO, msg)
}

func CheckPrint(e error) {
	if e != nil {
		//debug.PrintStack()
		Logger.Printlnf(Logging.ERROR, "%v", e)
	}
}

func Read(r io.Reader) ([]byte, error) {
	b := make([]byte, 1024*8)
	l, e := r.Read(b)
	return b[:l], e
}

func ExecuteCommand(cmmnd string) string {
	c := exec.Command("bash", "-c", cmmnd)
	b, e := c.Output()
	CheckPrint(e)

	return string(b)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Not secure random string generator, only for channel IDs or simple stuff
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
