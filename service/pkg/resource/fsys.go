package resource

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed fs/*
var fsys embed.FS

func ReadFile(name string) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("invalid argument : file name is empty")
	}
	//if fsys == nil {
	//	return nil, fmt.Errorf("invalid argument : file system is nil")
	//}
	return fs.ReadFile(fsys, name)
}

func ReadMap(name string) (map[string]string, error) {
	buf, err := ReadFile(name)
	if err != nil {
		return nil, err
	}
	return ParseBuffer(buf)

}
