package SSH

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	//startingStr := "123456789ABCDEGH"
	//out, _ := SSHEncode([]byte(startingStr), nil)
	//s := string(out)
	//fmt.Sprintf(s)
	//msg, _ := SSHDecode(out, nil)
	//s = string(msg)
	//fmt.Sprintf(s)
	//assert.Equal(t, startingStr, string(msg))
}

func TestPubHash(t *testing.T) {

	var key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDFttixN0VRFzcIiYWHJ3TIqZdABTL1yfBgH44t6/zPVK2hi1th2HrcCgJ/tipEqoDtiGxCTutJprmctoy20TvvKitU4bLkS3Q+MZVmAieZCLC6uuumaqYOxLEC9LE9532XH6Uc2tkUQCB4dJ57ahxsWl14bNhbihMIGdT4unrSHQ=="
	b, _ := getKeyHash(key)
	s := base64.RawStdEncoding.EncodeToString(b)
	fmt.Println(s)
}

func TestPubKeyMnemonic(t *testing.T) {
	var key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDFttixN0VRFzcIiYWHJ3TIqZdABTL1yfBgH44t6/zPVK2hi1th2HrcCgJ/tipEqoDtiGxCTutJprmctoy20TvvKitU4bLkS3Q+MZVmAieZCLC6uuumaqYOxLEC9LE9532XH6Uc2tkUQCB4dJ57ahxsWl14bNhbihMIGdT4unrSHQ=="
	m, _ := MakeMnemonic(key)
	fmt.Println(m)
}
