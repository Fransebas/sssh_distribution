/*
This is a wrapper for the ssh.channels so that I'm able to send packages of the size I want to.

What this does is break the packages in to chunks of  2^14 bytes, add them an ID and numerate every chunk

Currently only the server makes this, I'm not sure if the client should to or I should implement a more advance method for big chunks
of data in ssh, maybe something like http. Because with this method the whole data has to be stored in memory and can be properly buffered
*/
package LimitlessChannel

import (
	"encoding/json"
	"io"
	"rgui/CustomUtils"
)

// The writer most surely is a ssh.Channel but I leave it general
type LimitlessChannel struct {
	Writer io.Writer
}

type Package struct {
	Message       string
	Size          int
	TotalPackages int
	PackageNumber int
	ID            string
	mode          string // binary or string
}

// This is the size of the packages, the maximum set by ssh is 2^15 and to make sure the package fits I'm going to use 2^14
// const PACKAGE_SIZE  = 1 << 14
const PACKAGE_SIZE = 1 << 3

func (lc LimitlessChannel) Write(p []byte) (n int, err error) {
	l := len(p)
	ID := CustomUtils.RandStringRunes(20) // possibility of collision is 1/(52^20)
	for i := 0; i < l; i += PACKAGE_SIZE {
		// TODO: proper handle of error #1242345
		j := MinOf(len(p), i+PACKAGE_SIZE)
		packagePart := Package{
			Message:       string(p[i:j]),
			Size:          j - i,
			TotalPackages: (l + PACKAGE_SIZE - 1) / PACKAGE_SIZE, // round up
			PackageNumber: i / PACKAGE_SIZE,                      // round down
			mode:          "string",                              // For now all will be string
			ID:            ID,
		}

		b, _ := json.Marshal(packagePart)

		_, e := lc.Writer.Write(b)

		CustomUtils.CheckPrint(e)
	}
	return l, nil
}

func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
