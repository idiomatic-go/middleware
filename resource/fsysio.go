package resource

import (
	"fmt"
	"io/fs"
)

const (
	//token = byte(' ')
	eol       = byte('\n')
	comment   = "//"
	delimiter = ":"
	//root      = "resource"
)

func readFile(fsys fs.FS, name string) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("invalid argument : file name is empty")
	}
	if fsys == nil {
		return nil, fmt.Errorf("invalid argument : file system is nil")
	}
	return fs.ReadFile(fsys, name)
}

/*
func ReadFileContext(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid argument : context is nil")
	}
	return ReadFile(ContextEmbeddedFS(ctx), ContextEmbeddedContent(ctx))
}

*/

func ReadDir(fsys fs.FS, name string) ([]fs.DirEntry, error) {
	if fsys == nil {
		return nil, fmt.Errorf("invalid argument : file system is nil")
	}
	if name == "" {
		return nil, fmt.Errorf("invalid argument : directory name is empty")
	}
	return fs.ReadDir(fsys, name)
}

func readMap(fsys fs.FS, name string) (map[string]string, error) {
	buf, err := readFile(fsys, name)
	if err != nil {
		return nil, err
	}
	return ParseBuffer(buf)
}
