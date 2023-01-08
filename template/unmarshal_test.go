package template

import (
	"bytes"
	"fmt"
	"github.com/idiomatic-go/middleware/template/internal"
	"net/http"
	"strings"
)

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

func Example_Unmarshal() {

	result, status := Unmarshal[string](nil)
	fmt.Printf("test: Unmarshal(nil) -> [T:%v] [status:%v]\n", result, status)

	resp := new(http.Response)
	result, status = Unmarshal[string](resp)
	fmt.Printf("test: Unmarshal(resp) -> [T:%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: strings.NewReader("string: Hello World"), Err: nil}
	result, status = Unmarshal[string](resp)
	fmt.Printf("test: Unmarshal(resp) -> [T:%v] [status:%v]\n", result, status)

	resp.Body = &internal.ReaderCloser{Reader: bytes.NewReader([]byte("[]byte: Hello World")), Err: nil}
	result2, status2 := Unmarshal[[]byte](resp)
	fmt.Printf("test: Unmarshal(resp) -> [T:%v] [status:%v]\n", string(result2), status2)

	//Output:
	//test: Unmarshal(nil) -> [T:] [status:-1 Invalid Content [response or response body is nil]]
	//test: Unmarshal(resp) -> [T:] [status:-1 Invalid Content [response or response body is nil]]
	//test: Unmarshal(resp) -> [T:string: Hello World] [status:0 Successful]
	//test: Unmarshal(resp) -> [T:[]byte: Hello World] [status:0 Successful]

}
