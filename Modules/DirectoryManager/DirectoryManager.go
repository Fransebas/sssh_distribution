package DirectoryManager

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
	"sssh_server/CustomUtils"
)

type DirectoryManager struct {
	UserDirectory  string
	ConfigFolder   string
	VariableFolder string
}

func getConfigFolder(homeArg, pathArg, username string) string {
	home := path.Clean(homeArg)
	basePath := path.Clean(pathArg)
	dir := ""
	if runtime.GOOS == "windows" {
		// running no Windows
		panic("Windows not supported yet =(")
	}

	if runtime.GOOS == "darwin" {
		// running no MacOS
		dir = path.Join(home, "/Library/Application Support/sssh_server")
	} else {
		// linux and something else
		// TODO: validate other OS's
		dir = path.Join(home, "/.sssh_server")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory if not exist
		// It is created like this to avoid permissions errors
		CustomUtils.ExecuteCommandOnce(fmt.Sprintf(`mkdir "%v"`, dir), username)
	}

	return path.Clean(path.Join(dir, basePath))
}

func getVariableFolder(homeArg, pathArg, username string) string {
	home := path.Clean(homeArg)
	basePath := path.Clean(pathArg)
	dir := ""
	if runtime.GOOS == "windows" {
		// running no Windows
		panic("Windows not supported yet =(")
	}

	if runtime.GOOS == "darwin" {
		// running no MacOS
		dir = path.Join(home, "/Library/Application Support/sssh_server/var")
	} else {
		// linux and something else
		// TODO: validate other OS's
		dir = path.Join(home, "/.sssh_server")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory if not exist
		// It is created like this to avoid permissions errors
		CustomUtils.ExecuteCommandOnce(fmt.Sprintf(`mkdir "%v"`, dir), username)
	}

	return path.Clean(path.Join(dir, basePath))
}

func New(username string) *DirectoryManager {
	usr, _ := user.Lookup(username)
	s := usr.HomeDir
	dm := DirectoryManager{
		UserDirectory:  s,
		ConfigFolder:   getConfigFolder(s, "", username),
		VariableFolder: getVariableFolder(s, "", username),
	}
	return &dm
}

func (d *DirectoryManager) GetConfigFile(fileName string) string {
	return path.Join(d.ConfigFolder, fileName)
}

func (d *DirectoryManager) GetVariableFile(fileName string) string {
	return path.Join(d.VariableFolder, fileName)
}
