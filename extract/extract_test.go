package extract

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
	"time"
)

func Example_Initialize_Url() {
	err := Initialize("", nil, nil)
	fmt.Printf("test: initialize(\"\") -> [%v] [url:%v]\n", err, url)

	err = Initialize("test", nil, nil)
	fmt.Printf("test: initialize(\"\") -> [%v] [url:%v]\n", err, url)

	//Output:
	//test: initialize("") -> [invalid argument : uri is empty] [url:]
	//test: initialize("") -> [<nil>] [url:test]

}

func Example_Handler_NotProcessed() {
	url = "http://localhost:8080/accesslog"

	status := handler(nil)
	fmt.Printf("test: handler(nil) -> [%v]\n", status)

	req, _ := http.NewRequest("post", "http://localhost:8080/accesslog", nil)
	data := new(accesslog.Logd)
	data.AddRequest(req)
	status = handler(data)
	fmt.Printf("test: handler(data) -> [%v]\n", status)

	//Output:
	//test: handler(nil) -> [false]
	//test: handler(data) -> [false]

}

func Example_Handler_ConnectFailure() {
	url = "http://localhost:8080/accesslog"

	req, _ := http.NewRequest("post", "localhost:8081/accesslog", nil)
	data := new(accesslog.Logd)
	data.AddRequest(req)
	status := handler(data)
	fmt.Printf("test: handler(data) -> [%v]\n", status)

	//Output:
	//test: handler(data) -> [false]
}

func Example_Handler_Processed() {
	// Override the message handler
	handler = func(l *accesslog.Logd) bool {
		fmt.Printf("test: handler(logd) -> [%v]\n", accesslog.FormatJson(entries, l))
		return true
	}

	err := Initialize("http://localhost:8080/accesslog", nil, nil)
	fmt.Printf("test: initialize() -> [%v]\n", err)

	r0, _ := route.NewRoute("route-data-0")
	r1, _ := route.NewRoute("route-data-1")
	r2, _ := route.NewRoute("route-data-2")
	r3, _ := route.NewRoute("route-data-3")

	data0 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-1"}, Route: r0}
	data1 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-2"}, Route: r1}
	data2 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-3"}, Route: r2}
	data3 := accesslog.Logd{Origin: &accesslog.Origin{Region: "region-4"}, Route: r3}
	extract(&data0)
	extract(&data1)
	extract(&data2)
	extract(&data3)
	time.Sleep(time.Second * 2)
	Shutdown()

	//Output:
	//test: initialize() -> [<nil>]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":"route-data-0","traffic":null,"region":"region-1","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","response_flags":null,"bytes_received":"0","bytes_sent":"0","timeout":-1,"limit":"INF","burst":1}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":"route-data-1","traffic":null,"region":"region-2","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","response_flags":null,"bytes_received":"0","bytes_sent":"0","timeout":-1,"limit":"INF","burst":1}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":"route-data-2","traffic":null,"region":"region-3","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","response_flags":null,"bytes_received":"0","bytes_sent":"0","timeout":-1,"limit":"INF","burst":1}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"route_name":"route-data-3","traffic":null,"region":"region-4","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","response_flags":null,"bytes_received":"0","bytes_sent":"0","timeout":-1,"limit":"INF","burst":1}]

}
