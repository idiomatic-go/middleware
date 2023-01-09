package vhost

import (
	"embed"
	"fmt"
)

//go:embed resource/*
var content embed.FS

func init() {
	MountFS(content)
}

func _ExampleFileSystemNotMounted() {
	_, err := ReadFile("resource/readme.txt")
	fmt.Printf("Error : %v\n", err)

	//Output:
	// Error : invalid argument : file system has not been mounted
}

func ExampleReadFile() {
	_, err0 := ReadFile("")
	fmt.Printf("Error : %v\n", err0)

	buf, err := ReadFile("bad-path/config_bad.txt")
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	} else {
		fmt.Println(string(buf))
	}

	buf, err = ReadFile("postgresql/config_{env}.txt")
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	} else {
		fmt.Println(string(buf))
	}

	// Should override and return config_test.txt
	/*
		lookupEnv = func(name string) (string, error) { return "stage", nil }
		buf, err = ReadFile("postgresql/config_{env}.txt")
		if err != nil {
			fmt.Printf("Error : %v\n", err)
		} else {
			fmt.Println(string(buf))
		}
	*/

	//Output:
	// Error : invalid argument : file name is empty
	// Error : open resource/bad-path/config_bad.txt: file does not exist
	// Error : invalid argument : template variable is invalid: env
}
