package template

import "fmt"

func ExampleDebugHandler_Handle() {
	var h DebugHandler

	fmt.Printf("test: Handle() -> %v\n", h.Handle("location", nil))

	//Output:
	//test: Handle() -> 0 The operation was successful []
}
