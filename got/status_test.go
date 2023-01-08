package got

import (
	"errors"
	"fmt"
	"net/http"
)

func ExampleStatus_String() {
	s := NewStatus(StatusPermissionDenied, "", nil)
	fmt.Printf("test: NewStatus() -> [%v]\n", s)

	s = NewStatus(StatusOutOfRange, "", errors.New("error - 1"), errors.New("error - 2"))
	fmt.Printf("test: NewStatus() -> [%v]\n", s)

	//Output:
	//test: NewStatus() -> [0 Successful]
	//test: NewStatus() -> [11 The operation was attempted past the valid range [error - 1 error - 2]]
}

func ExampleStatus_Http() {
	location := "test"
	err := errors.New("http error")
	fmt.Printf("test: NewHttpStatus(nil) -> [%v]\n", NewHttpStatus(nil, location, nil))
	fmt.Printf("test: NewHttpStatus(nil) -> [%v]\n", NewHttpStatus(nil, location, err))

	resp := http.Response{StatusCode: http.StatusBadRequest}
	fmt.Printf("test: NewHttpStatus(resp) -> [%v]\n", NewHttpStatus(&resp, location, nil))
	fmt.Printf("test: NewHttpStatus(resp) -> [%v]\n", NewHttpStatus(&resp, location, err))

	//Output:
	//test: NewHttpStatus(nil) -> [-1 Invalid Content]
	//test: NewHttpStatus(nil) -> [500 Internal Error [http error]]
	//test: NewHttpStatus(resp) -> [400 Bad Request]
	//test: NewHttpStatus(resp) -> [500 Internal Error [http error]]

}
