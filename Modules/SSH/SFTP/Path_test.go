package SFTP

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestPaths(t *testing.T) {
	path := "/Users/fransebas/CLionProjects/roadmap/CF-B/UVA10325/cmake-build-debug"
	r := filepath.Clean(path)
	paths := filepath.SplitList(path)
	filepath.Join()
	filepath.Dir(path)
	r2 := filepath.FromSlash(path)
	fmt.Println(paths)
	fmt.Println(r)
	fmt.Println(r2)
}
