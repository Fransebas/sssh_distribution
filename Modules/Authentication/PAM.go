/*
This module manage the authentication, it's based on the PAM standard.
It should work on most distros but I'm not sure about using the service `sshd`

PAM documentation

	* http://www.linux-pam.org/Linux-PAM-html/Linux-PAM_ADG.html
	* https://likegeeks.com/linux-pam-easy-guide/#Linux-PAM-Configuration
	* Files located at /etc/pam.d
	* maybe? https://github.com/uber/pam-ussh
*/

package Authentication

//#cgo LDFLAGS: -lpam
/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/Fransebas/golang-pam"

	//"github.com/vvanpo/golang-pam"
	"sssh_server/CustomUtils"
	"strconv"
	"strings"
	"unsafe"
)

// Returns true if the given username exist in the system, it simply uses id -u 'username'
func UserExist(username string) bool {
	res := CustomUtils.ExecuteCommand(fmt.Sprintf("id -u %v", username))
	if len(res) <= 0 {
		return false
	}
	res = res[:len(res)-1]
	parts := strings.Split(res, " ")

	// Double check
	if len(parts) > 1 {
		return false
	}
	_, e := strconv.Atoi(res)

	if e != nil {
		return false
	}
	return true
}

// Returns true if the username and password are correct
func ValidateUser(username, password string) (bool, error) {
	// uses the PAM library to validate user and password

	// For some reason I can only send one string, that's why I need to parse to json the object
	UserPssword := userPassword{
		Username: username,
		Password: password,
	}
	b, _ := json.Marshal(UserPssword)

	a := C.CString(string(b))
	conn := myConHand{
		username: a,
	}

	// sshd service should work on most linux distros but I don't know when it doesn't
	tx, _ := pam.Start("sshd", username, conn)
	r := tx.Authenticate(0)

	C.free(unsafe.Pointer(a))
	tx.End(0)
	if r == pam.SUCCESS {
		return true, nil
	} else {
		return false, nil
	}

}

type userPassword struct {
	Username string
	Password string
}

type myConHand struct {
	username *C.char
}

//C.GoString
func (m myConHand) RespondPAM(msg_style int, msg string) (string, bool) {
	str := C.GoString(m.username)

	var UserPssword userPassword
	_ = json.Unmarshal([]byte(str), &UserPssword)

	switch msg_style {
	case pam.PROMPT_ECHO_OFF:
		// Here pam asks for the password
		return UserPssword.Password, true
	case pam.PROMPT_ECHO_ON:
		// Here pam asks for the username
		return UserPssword.Username, true
	case pam.ERROR_MSG:
		fmt.Println(msg)
		return "", true
	case pam.TEXT_INFO:
		fmt.Println(msg)
		return "", true
	default:
		return "", true
	}
}
