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
	fmt.Printf("Empty : %v\n", IsEmpty(s))

	s = "    "
	fmt.Printf("Empty : %v\n", IsEmpty(s))

	s = "   def45 "
	fmt.Printf("Empty : %v\n", IsEmpty(s))

	//Output:
	//Empty : true
	//Empty : true
	//Empty : false
}
