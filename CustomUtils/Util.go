package CustomUtils

import "fmt"

func CheckPanic(e error, msg string) {
	if e != nil {
		panic(e.Error() + msg)
	}
}

func CheckPrint(e error) {
	if e != nil {
		fmt.Println("error : " + e.Error())
	}
}
