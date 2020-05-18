package SFTP

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

const (
	Setstat = "Setstat"
	Rename  = "Rename"
	Rmdir   = "Rmdir"
	Mkdir   = "Mkdir"
	Link    = "Link"
	Symlink = "Symlink"
	Remove  = "Remove"
)

const (
	List     = "List"
	Stat     = "Stat"
	Readlink = "Readlink"
)

type SFTPRequestServer struct {
	user string
}

func New(user string) *SFTPRequestServer {
	h := new(SFTPRequestServer)
	h.user = user
	return h
}

func (s *SFTPRequestServer) Filecmd(r *sftp.Request) error {
	path := r.Filepath
	target := r.Target

	switch r.Method {
	case Setstat:
		if !s.CanWrite(path) {
			return sftp.ErrSSHFxPermissionDenied
		}
		r.Attributes()
		_, e := os.Stat(path)
		return e
	case Rename:
		if !s.CanWrite(path) {
			return sftp.ErrSSHFxPermissionDenied
		}
		return os.Rename(path, target)
	case Rmdir, Remove:
		if !s.CanWrite(path) {
			return sftp.ErrSSHFxPermissionDenied
		}
		return os.Remove(path)
	case Mkdir:
		// TODO: is this the best permission for a new dir?
		pPath := prevPath(path)
		// Can create in the upper directory?
		if !s.CanWrite(pPath) {
			return sftp.ErrSSHFxPermissionDenied
		}
		e := os.Mkdir(path, 0777)
		if e != nil {
			return e
		}
		// Here we change the permissions of the directory because its root
		// At this point we don't care about the errors because the dir was made
		_ = chown(s.user, path)
		return nil
	case Link:
		// TODO: should I change permissions of link?
		if !s.CanWrite(path) {
			return sftp.ErrSSHFxPermissionDenied
		}

		pPath := prevPath(target)
		// Can create in the folder is going to be in?
		if !s.CanWrite(pPath) {
			return sftp.ErrSSHFxPermissionDenied
		}
		return os.Link(path, target)
	case Symlink:
		// TODO: should I change permissions of symblink?
		if !s.CanWrite(path) {
			return sftp.ErrSSHFxPermissionDenied
		}
		return os.Symlink(path, target)
	default:
		return errors.New(fmt.Sprintf("Command %v not supported", r.Method))
	}
}

// change path owner to username
func chown(username, path string) (e error) {
	usr, e := user.Lookup(username)
	if e != nil {
		return
	}
	gid, e := strconv.Atoi(usr.Gid)
	if e != nil {
		return
	}
	uid, e := strconv.Atoi(usr.Uid)
	if e != nil {
		return
	}
	e = os.Chown(path, uid, gid)
	return
}

func (s *SFTPRequestServer) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	path := r.Filepath
	if s.CanRead(path) {
		return os.Open(path)
	}
	return nil, sftp.ErrSSHFxPermissionDenied
}

func (s *SFTPRequestServer) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	path := r.Filepath

	if _, err := os.Stat(path); err == nil {
		// path exists
		if s.CanWrite(path) {
			return os.Create(path)
		}
	} else if os.IsNotExist(err) {
		// path does *not* exist
		// check if we can write it that given folder
		pPath := prevPath(path)
		if !s.CanWrite(pPath) {
			return nil, sftp.ErrSSHFxPermissionDenied
		}

		f, e := os.Create(path)
		if e != nil {
			return nil, e
		}
		// ignore error because at this point the file exist so is futile =(
		e = chown(s.user, path)
		return f, e
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return nil, err
	}

	return nil, sftp.ErrSshFxPermissionDenied
}

func (s *SFTPRequestServer) Filelist(r *sftp.Request) (sftp.ListerAt, error) {

	path := r.Filepath
	if !(s.CanExecute(path) || s.CanRead(path)) {
		//fmt.Printf("ERROR File list method %v in path %v\n", r.Method, path)
		return nil, sftp.ErrSSHFxPermissionDenied
	}

	switch r.Method {
	case List:
		f, e := os.Open(r.Filepath)
		defer f.Close()
		if e != nil {
			return nil, e
		}
		l, e := f.Readdir(-1)
		if e != nil {
			return nil, e
		}
		return listAt(l), nil
	case Stat:
		fi, e := os.Stat(path)
		return listAt([]os.FileInfo{fi}), e
	case Readlink:
		return nil, errors.New("Readlink not implemented yet =(")
	default:
		return nil, errors.New(fmt.Sprintf("Method %v not supported", r.Method))
	}
}

func (s *SFTPRequestServer) testPermission(path, mode string) bool {
	cmmnd := fmt.Sprintf(`sudo -u %v test -%v "%v"; echo "$?"`, s.user, mode, path)
	c := exec.Command("bash", "-c", cmmnd)
	b, e := c.Output()

	if e != nil {
		return false
	}

	return string(b)[:1] == "0"
}

// Returns the path one step up the hierarchy
// the path should be absolute without any ../ or ./
func prevPath(path string) string {
	// TODO : this functinos is not 100% correct, we need to se if there is no way of having the character /
	// in the path string like \/ or something
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "/")
	}

	// if len(parts) <= 1 that means that it was the root directory
	// and it should return /
	return "/"

}

func (s *SFTPRequestServer) CanWrite(path string) bool {
	return s.testPermission(path, "w")
}

func (s *SFTPRequestServer) CanRead(path string) bool {
	return s.testPermission(path, "r")
}

func (s *SFTPRequestServer) CanExecute(path string) bool {
	return s.testPermission(path, "x")
}

type listAt []os.FileInfo

// ListAt returns the number of entries copied and an io.EOF error if we made it to the end of the file list.
// Take a look at the pkg/sftp godoc for more information about how this function should work.
func (l listAt) ListAt(f []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}

	n := copy(f, l[offset:])
	if n < len(f) {
		return n, io.EOF
	}
	return n, nil
}
