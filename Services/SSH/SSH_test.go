package SSH

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryption(t *testing.T) {
	startingStr := "123456789ABCDEGH"
	out, _ := SSHEncode([]byte(startingStr), nil)
	s := string(out)
	fmt.Sprintf(s)
	msg, _ := SSHDecode(out, nil)
	s = string(msg)
	fmt.Sprintf(s)
	assert.Equal(t, startingStr, string(msg))
}
