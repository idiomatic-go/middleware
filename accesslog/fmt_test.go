package accesslog

import "fmt"

func NilEmpty(s string) string {
	if IsEmpty(s) {
		return "<nil>"
	}
	return s
}

func ExampleIsEmpty() {
	var s = ""
	fmt.Printf("test: Empty() -> %v\n", IsEmpty(s))

	s = "    "
	fmt.Printf("test: Empty() -> %v\n", IsEmpty(s))

	s = "   def45 "
	fmt.Printf("test: Empty() -> %v\n", IsEmpty(s))

	//Output:
	//test: Empty() -> true
	//test: Empty() -> true
	//test: Empty() -> false
}
