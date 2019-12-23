package SSH

import (
	"fmt"
	"io"
	"testing"
)

func TestSSSHServer_ListenAndServe(t *testing.T) {
	server := SSSHServer{}

	server.HandleFunc("echo", func(s *SSHSession, w io.Writer, r io.Reader) {
		io.Copy(w, r)
	})

	server.HandleFunc("sessionid", func(s *SSHSession, w io.Writer, r io.Reader) {
		w.Write([]byte(s.GetSessionID()))
		fmt.Printf("session %v \n", s.GetSessionID())
	})

	server.ListenAndServe()
}
