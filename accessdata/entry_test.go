package accessdata

import (
	"fmt"
	"net/http"
)

func Example_Value_Origin() {
	op := OriginRegionOperator
	o := Origin{"region", "zone", "", "", ""}
	data := Entry{Origin: nil}
	fmt.Printf("test: Value(\"region\") -> [%v]\n", data.Value(op))

	data = Entry{Origin: &o}
	fmt.Printf("test: Value(\"region\") -> [%v]\n", data.Value(op))

	//Output:
	//test: Value("region") -> []
	//test: Value("region") -> [region]
}

func Example_Value_Actuator() {
	name := "test-route"
	op := RouteNameOperator

	data := Entry{}
	fmt.Printf("test: Value(\"%v\") -> [%v]\n", name, data.Value(op))

	data = Entry{ActState: map[string]string{ActName: name}}
	fmt.Printf("test: Value(\"%v\") -> [route_name:%v]\n", name, data.Value(op))

	data = Entry{ActState: map[string]string{TimeoutName: "500"}}
	fmt.Printf("test: Value(\"%v\") -> [timeout:%v]\n", name, data.Value(TimeoutDurationOperator))

	//Output:
	//test: Value("test-route") -> []
	//test: Value("test-route") -> [route_name:test-route]
	//test: Value("test-route") -> [timeout:500]
}

func Example_Value_Request() {
	op := RequestMethodOperator

	data := &Entry{}
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(op))

	req, _ := http.NewRequest("POST", "www.google.com", nil)
	//req.Header.Add(RequestIdHeaderName, uuid.New().String())
	req.Header.Add(RequestIdHeaderName, "123-456-789")
	req.Header.Add(FromRouteHeaderName, "calling-route")
	data = &Entry{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(op))

	fmt.Printf("test: Value(\"headers\") -> [request-id:%v] [from-route:%v]\n", data.Value(RequestIdOperator), data.Value(RequestFromRouteOperator))

	//Output:
	//test: Value("method") -> []
	//test: Value("method") -> [POST]
	//test: Value("headers") -> [request-id:123-456-789] [from-route:calling-route]
}

func Example_Value_Response() {
	op := ResponseStatusCodeOperator

	data := &Entry{}
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(op))

	resp := &http.Response{StatusCode: 200}
	data = &Entry{}
	data.AddResponse(resp)
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(op))

	//Output:
	//test: Value("code") -> [0]
	//test: Value("code") -> [200]
}

func Example_Value_Request_Header() {
	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	data := Entry{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"customer\") -> [%v]\n", data.Value("%REQ(customer)%"))

	//Output:
	//test: Value("customer") -> [Ted's Bait & Tackle]
}
