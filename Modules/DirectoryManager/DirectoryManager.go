package DirectoryManager

import (
	"fmt"
	"os/user"
	"runtime"
	"sssh_server/CustomUtils"
)

type DirectoryManager struct {
	UserDirectory  string
	ConfigFolder   string
	VariableFolder string
}

func getConfigFolder(homeArg, pathArg, username string) string {
	home := removeEndSlash(homeArg)
	path := removeStartSlash(pathArg)
	dir := ""
	if runtime.GOOS == "windows" {
		// running no Windows
		panic("Windows not supported yet =(")
	}

	if runtime.GOOS == "darwin" {
		// running no MacOS
		dir = fmt.Sprintf("%v/Library/Application Support/sssh_server", home)
	} else {
		// linux and something else
		// TODO: validate other OS's
		dir = fmt.Sprintf("%v/.sssh_server", home)
	}

	// Create directory if not exist
	// It is created like this to avoid permissions errors
	CustomUtils.ExecuteCommand(fmt.Sprintf(`sudo -u %v mkdir "%v"`, username, dir))

	return removeEndSlash(fmt.Sprintf("%v/%v", dir, path))
}

func getVariableFolder(homeArg, pathArg, username string) string {
	home := removeEndSlash(homeArg)
	path := removeStartSlash(pathArg)
	dir := ""
	if runtime.GOOS == "windows" {
		// running no Windows
		panic("Windows not supported yet =(")
	}

	if runtime.GOOS == "darwin" {
		// running no MacOS
		dir = fmt.Sprintf("%v/Library/Application Support/sssh_server/var", home)
	} else {
		// linux and something else
		// TODO: validate other OS's
		dir = fmt.Sprintf("%v/.sssh_server", home)
	}

	// Create directory if not exist
	// It is created like this to avoid permissions errors
	CustomUtils.ExecuteCommand(fmt.Sprintf(`sudo -u %v mkdir "%v"`, username, dir))

	return removeEndSlash(fmt.Sprintf("%v/%v", dir, path))
}

func removeEndSlash(path string) string {
	if len(path) <= 0 {
		return path
	}
	if path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	return path
}

func removeStartSlash(path string) string {
	if len(path) <= 0 {
		return path
	}
	if path[0] == '/' {
		return path[1:]
	}
	return path
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
	return removeEndSlash(fmt.Sprintf("%v/%v", removeEndSlash(d.ConfigFolder), removeStartSlash(fileName)))
}

func (d *DirectoryManager) GetVariableFile(fileName string) string {
	return removeEndSlash(fmt.Sprintf("%v/%v", removeEndSlash(d.VariableFolder), removeStartSlash(fileName)))
}
