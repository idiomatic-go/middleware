package resource

import (
	"embed"
	"fmt"
)

//go:embed resource/*
var content embed.FS

func init() {
	MountFS(content)
}

func ExampleFileSystemNotMounted() {
	name := "resource/readme.txt"
	_, err := ReadFile(name)
	fmt.Printf("test: ReadFile(%v) -> [error:%v]\n", name, err)

	//Output:
	//test: ReadFile(resource/readme.txt) -> [error:open resource/resource/readme.txt: file does not exist]

}

func _ExampleReadFile() {
	_, err0 := ReadFile("")
	fmt.Printf("test: ReadFile() -> %v\n", err0)

	name := "bad-path/config_bad.txt"
	buf, err := ReadFile(name)
	fmt.Printf("test: ReadFile(%v) -> [error:%v] [content:%v]\n", name, err, string(buf))

	name = "postgresql/config_{env}.txt"
	buf, err = ReadFile(name)
	fmt.Printf("test: ReadFile(%v) -> [error:%v] [content:%v]\n", name, err, string(buf))

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
	//test: ReadFile() -> invalid argument : file name is empty
	//test: ReadFile(bad-path/config_bad.txt) -> [error:open resource/bad-path/config_bad.txt: file does not exist] [content:]
	//test: ReadFile(postgresql/config_{env}.txt) -> [error:<nil>] [content:// this is the test environment
	//env : dev
	//next  : second value
	//timeout : 10020]

}
