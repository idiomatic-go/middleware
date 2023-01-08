package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/idiomatic-go/middleware/template/internal"
	"net/http"
	"strings"
)

type adddressIO struct {
	Name   string
	Street string
	City   string
	State  string
}

/*
func _ExampleString_Unmarshal() {
	u := new(StringUnmarshal)

	result, status := UnmarshalInterface[string, DebugHandler](nil, u)
	fmt.Printf("test: Unmarshal(nil) -> [result:%v] [status:%v]\n", result, status)

	resp := new(http.Response)
	result, status = UnmarshalInterface[string, DebugHandler](resp, u)
	fmt.Printf("test: Unmarshal(resp) -> [result:%v] [status:%v]\n", result, status)

	//Output:
	//fail
}

*/

func Example_Unmarshal() {
	result, status := Unmarshal[string](nil)
	fmt.Printf("test: Unmarshal[string](nil) -> [%v] [status:%v]\n", result, status)

	resp := new(http.Response)
	result, status = Unmarshal[string](resp)
	fmt.Printf("test: Unmarshal[string](resp) -> [%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: strings.NewReader("Hello World String"), Err: nil}
	result, status = Unmarshal[string](resp)
	fmt.Printf("test: Unmarshal[string](resp) -> [%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: bytes.NewReader([]byte("Hello World []byte")), Err: nil}
	result2, status2 := Unmarshal[[]byte](resp)
	fmt.Printf("test: Unmarshal[[]byte](resp) -> [%v] [status:%v]\n", string(result2), status2)

	//Output:
	//test: Unmarshal[string](nil) -> [] [status:-1 Invalid Content [response or response body is nil]]
	//test: Unmarshal[string](resp) -> [] [status:-1 Invalid Content [response or response body is nil]]
	//test: Unmarshal[string](resp) -> [Hello World String] [status:0 Successful]
	//test: Unmarshal[[]byte](resp) -> [Hello World []byte] [status:0 Successful]

}

func Example_Unmarshal_Struct() {
	addr := adddressIO{
		Name:   "Bob Smith",
		Street: "123 Oak Avenue",
		City:   "New Orleans",
		State:  "LA",
	}

	buf, _ := json.Marshal(&addr)
	resp := new(http.Response)
	resp.Body = &internal.ReaderCloser{Reader: bytes.NewReader(buf), Err: nil}

	result, status := Unmarshal[adddressIO](resp)
	fmt.Printf("test: Unmarshal(resp) -> [T:%v] [status:%v]\n", result, status)

	//result, status = Decode[adddressIO](resp)
	//fmt.Printf("test: Unmarshal(resp) -> [T:%v] [status:%v]\n", result, status)

	//Output:
	//test: Unmarshal(resp) -> [T:{Bob Smith 123 Oak Avenue New Orleans LA}] [status:0 Successful]

}

func Example_Decode() {
	result, status := Decode[string](nil)
	fmt.Printf("test: Decode[string](nil) -> [%v] [status:%v]\n", result, status)

	resp := new(http.Response)
	result, status = Decode[string](resp)
	fmt.Printf("test: Decode[string](resp) -> [%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: strings.NewReader("Hello World String"), Err: nil}
	result, status = Decode[string](resp)
	fmt.Printf("test: Decode[string](resp) -> [%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: bytes.NewReader([]byte("Hello World []byte")), Err: nil}
	result2, status2 := Decode[[]byte](resp)
	fmt.Printf("test: Decode[[]byte](resp) -> [%v] [status:%v]\n", string(result2), status2)

	//Output:
	//test: Decode[string](nil) -> [] [status:-1 Invalid Content [response or response body is nil]]
	//test: Decode[string](resp) -> [] [status:-1 Invalid Content [response or response body is nil]]
	//test: Decode[string](resp) -> [Hello World String] [status:0 Successful]
	//test: Decode[[]byte](resp) -> [Hello World []byte] [status:0 Successful]

}
