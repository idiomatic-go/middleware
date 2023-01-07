package template

import (
	"errors"
	"fmt"
)

func ExampleStatus_String() {
	s := NewStatus(StatusPermissionDenied, "", nil)
	fmt.Printf("test: NewStatus() -> [%v]\n", s)

	s = NewStatus(StatusOutOfRange, "", errors.New("error - 1"), errors.New("error - 2"))
	fmt.Printf("test: NewStatus() -> [%v]\n", s)

	//Output:
	//test: NewStatus() -> [7 The caller does not have permission to execute the specified operation []]
	//test: NewStatus() -> [11 The operation was attempted past the valid range [error - 1 error - 2]]
}