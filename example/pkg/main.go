package main

import (
	"fmt"
	"runtime"
)

func main() {
	displayRuntime()
}

func displayRuntime() {
	fmt.Println(fmt.Sprintf("vers : %v", runtime.Version()))
	fmt.Println(fmt.Sprintf("os   : %v", runtime.GOOS))
	fmt.Println(fmt.Sprintf("arch : %v", runtime.GOARCH))
	fmt.Println(fmt.Sprintf("cpu  : %v", runtime.NumCPU()))
}
