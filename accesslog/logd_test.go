package accesslog

import (
	"fmt"
	"net/http"
)

func Example_Value_Origin() {
	op := OriginRegionOperator
	o := Origin{"region", "zone", "", "", ""}
	data := Logd{Origin: nil}
	fmt.Printf("test: Value(\"region\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	data = Logd{Origin: &o}
	fmt.Printf("test: Value(\"region\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	//Output:
	//test: Value("region") -> []
	//test: Value("region") -> [region]
}

func Example_Value_Actuator() {
	name := "test-route"
	op := RouteNameOperator

	data := Logd{}
	fmt.Printf("test: Value(\"%v\") -> [%v]\n", name, data.Value(NewEntry(op, "", "", true)))

	data = Logd{ActState: map[string]string{ActName: name}}
	fmt.Printf("test: Value(\"%v\") -> [route_name:%v]\n", name, data.Value(NewEntry(op, "", "", true)))

	data = Logd{ActState: map[string]string{TimeoutName: "500"}}
	//fmt.Printf("test: Value(\"%v\") -> [route_name:%v]\n", name, data.Value(NewEntry(op, "", "", true)))
	fmt.Printf("test: Value(\"%v\") -> [timeout:%v]\n", name, data.Value(NewEntry(TimeoutDurationOperator, "", "", false)))

	//Output:
	//test: Value("test-route") -> []
	//test: Value("test-route") -> [route_name:test-route]
	//test: Value("test-route") -> [timeout:500]
}

func Example_Value_Request() {
	op := RequestMethodOperator

	data := &Logd{}
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	req, _ := http.NewRequest("POST", "www.google.com", nil)
	data = &Logd{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"method\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	//Output:
	//test: Value("method") -> []
	//test: Value("method") -> [POST]
}

func Example_Value_Response() {
	op := ResponseStatusCodeOperator

	data := &Logd{}
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	resp := &http.Response{StatusCode: 200}
	data = &Logd{}
	data.AddResponse(resp)
	fmt.Printf("test: Value(\"code\") -> [%v]\n", data.Value(NewEntry(op, "", "", true)))

	//Output:
	//test: Value("code") -> [0]
	//test: Value("code") -> [200]
}

func Example_Value_Request_Header() {
	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	data := Logd{}
	data.AddRequest(req)
	fmt.Printf("test: Value(\"customer\") -> [%v]\n", data.Value(NewEntry("header:customer", "", "", true)))

	//Output:
	//test: Value("customer") -> [Ted's Bait & Tackle]
}
