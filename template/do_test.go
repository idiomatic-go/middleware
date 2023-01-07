package template

import (
	"fmt"
	"io"
	"net/http"
)

func ExampleDo_InvalidArgument() {
	_, s := Do[NoOpHandler](nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	req, _ := http.NewRequest("", "http://www.google.com", nil)
	_, s = DoClient[DebugHandler](req, nil)
	fmt.Printf("test: DoClient(req,nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [13 Internal Error [invalid argument: request is nil]]
	//[github.com/idiomatic-go/middleware/template/Do [invalid argument: client is nil]]
	//test: DoClient(req,nil) -> [13 Internal Error]

}

func ExampleDo_HttpError() {
	req, _ := http.NewRequest(http.MethodGet, "echo://www.somestupidname.com?httpError=true", nil)
	resp, status := Do[DebugHandler](req)
	fmt.Printf("test: Do(req) -> [%v] [response:%v]\n", status, resp)

	//Output:
	//[github.com/idiomatic-go/middleware/template/Do [http: connection has been hijacked]]
	//test: Do(req) -> [500 Internal Error] [response:<nil>]

}

func ExampleDo_IOError() {
	req, _ := http.NewRequest(http.MethodGet, "echo://www.somestupidname.com?ioError=true", nil)
	resp, s := Do[DebugHandler](req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [body:%v]\n", s, resp != nil, resp.StatusCode, resp.Body != nil)

	defer resp.Body.Close()
	buf, s2 := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) : [%v] [body:%v]\n", s2, string(buf))

	//Output:
	//test: Do(req) -> [200 Successful] [resp:true] [statusCode:200] [body:true]
	//test: io.ReadAll(resp.Body) : [unexpected EOF] [body:]

}

func ExampleDo_Success() {
	var uri = "echo://www.somestupidname.com"
	uri += "?content-type=text/html"
	uri += "&content-length=1234"
	uri += "&body=<html><body><h1>Hello, World</h1></body></html>"
	req, err0 := http.NewRequest(http.MethodGet, uri, nil)
	if err0 != nil {
		fmt.Println("test: init() -> failure")
	}
	resp, status := Do[DebugHandler](req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [content-type:%v] [content-length:%v] [body:%v]\n",
		status, resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Header.Get("content-length"), resp.Body != nil)

	defer resp.Body.Close()
	buf, ioError := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) : [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do(req) -> [200 Successful] [resp:true] [statusCode:200] [content-type:text/html] [content-length:1234] [body:true]
	//test: io.ReadAll(resp.Body) : [err:<nil>] [body:<html><body><h1>Hello, World</h1></body></html>]

}
