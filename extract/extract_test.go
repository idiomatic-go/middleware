package extract

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
	"time"
)

func _ExampleInitializeUrl() {
	err := Initialize("", nil, nil)
	fmt.Printf("Error  : %v\n", err)
	fmt.Printf("Url    : %v\n", url)

	err = Initialize("test", nil, nil)
	fmt.Printf("Error  : %v\n", err)
	fmt.Printf("Url    : %v\n", url)

	//Output:
	//Error  : invalid argument : uri is empty
	//Url    :
	//Error  : <nil>
	//Url    : test
}

func _ExampleMessageHandlerNotProcessed() {
	url = "http://localhost:8080/accesslog"

	status := handler(nil)
	fmt.Printf("Status : %v\n", status)

	req, _ := http.NewRequest("post", "http://localhost:8080/accesslog", nil)
	data := new(accesslog.Logd)
	data.AddRequest(req)
	status = handler(data)
	fmt.Printf("Status : %v\n", status)

	//Output:
	//Status : false
	//Status : false

}

func _ExampleMessageHandlerConnectFailure() {
	url = "http://localhost:8080/accesslog"

	req, _ := http.NewRequest("post", "localhost:8081/accesslog", nil)
	data := new(accesslog.Logd)
	data.AddRequest(req)
	status := handler(data)
	fmt.Printf("Status : %v\n", status)

	//Output:
	//Status : false
}

func ExampleMessageHandlerProcessed() {
	// Override the message handler
	handler = func(l *accesslog.Logd) bool {
		fmt.Printf("%v\n", accesslog.FormatJson(entries, l))
		return true
	}

	err := Initialize("http://localhost:8080/accesslog", nil, nil)
	fmt.Printf("Error : %v\n", err)

	data0 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-1"}}
	data1 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-2"}}
	data2 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-3"}}
	data3 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-4"}}
	extract(&data0)
	extract(&data1)
	extract(&data2)
	extract(&data3)
	time.Sleep(time.Second * 2)
	Shutdown()

	//Output:
	//Error : <nil>
	//{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":null,"traffic":null,"region":"region-1","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"authority":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":null,"response_flags":null,"bytes_received":null}
	//{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":null,"traffic":null,"region":"region-2","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"authority":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":null,"response_flags":null,"bytes_received":null}
	//{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":null,"traffic":null,"region":"region-3","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"authority":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":null,"response_flags":null,"bytes_received":null}
	//{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":null,"traffic":null,"region":"region-4","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"authority":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":null,"response_flags":null,"bytes_received":null}
}
