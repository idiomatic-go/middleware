package extract

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
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
	data := new(accessdata.Entry)
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
	data := new(accessdata.Entry)
	data.AddRequest(req)
	status := handler(data)
	fmt.Printf("test: handler(data) -> [%v]\n", status)

	//Output:
	//test: handler(data) -> [false]
}

func Example_Handler_Processed() {
	// Override the message handler
	handler = func(l *accessdata.Entry) bool {
		fmt.Printf("test: handler(logd) -> [%v]\n", accessdata.WriteJson(entries, l))
		return true
	}

	err := Initialize("http://localhost:8080/accesslog", nil, nil)
	fmt.Printf("test: initialize() -> [%v]\n", err)

	//r0, _ := route.NewRoute("route-data-0")
	//r1, _ := route.NewRoute("route-data-1")
	//r2, _ := route.NewRoute("route-data-2")
	//r3, _ := route.NewRoute("route-data-3")

	data0 := accessdata.Entry{Origin: &accessdata.Origin{Region: "region-1"}, ActState: map[string]string{accessdata.ActName: "route-data-0"}}
	data1 := accessdata.Entry{Origin: &accessdata.Origin{Region: "region-2"}, ActState: map[string]string{accessdata.ActName: "route-data-1"}}
	data2 := accessdata.Entry{Origin: &accessdata.Origin{Region: "region-3"}, ActState: map[string]string{accessdata.ActName: "route-data-2"}}
	data3 := accessdata.Entry{Origin: &accessdata.Origin{Region: "region-4"}, ActState: map[string]string{accessdata.ActName: "route-data-3"}}
	extract(&data0)
	extract(&data1)
	extract(&data2)
	extract(&data3)
	time.Sleep(time.Second * 2)
	Shutdown()

	//Output:
	//test: initialize() -> [<nil>]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":null,"route_name":"route-data-0","region":"region-1","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","status_flags":null,"bytes_received":"0","bytes_sent":"0","timeout_ms":null,"rate_limit":null,"rate_burst":null,"retry":null,"retry_rate_limit":null,"retry_rate_burst":null,"failover":null}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":null,"route_name":"route-data-1","region":"region-2","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","status_flags":null,"bytes_received":"0","bytes_sent":"0","timeout_ms":null,"rate_limit":null,"rate_burst":null,"retry":null,"retry_rate_limit":null,"retry_rate_burst":null,"failover":null}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":null,"route_name":"route-data-2","region":"region-3","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","status_flags":null,"bytes_received":"0","bytes_sent":"0","timeout_ms":null,"rate_limit":null,"rate_burst":null,"retry":null,"retry_rate_limit":null,"retry_rate_burst":null,"failover":null}]
	//test: handler(logd) -> [{"start_time":"0001-01-01 00:00:00.000000","duration_ms":0,"traffic":null,"route_name":"route-data-3","region":"region-4","zone":null,"sub_zone":null,"service":null,"instance_id":null,"method":null,"host":null,"path":null,"protocol":null,"request_id":null,"forwarded":null,"status_code":"0","status_flags":null,"bytes_received":"0","bytes_sent":"0","timeout_ms":null,"rate_limit":null,"rate_burst":null,"retry":null,"retry_rate_limit":null,"retry_rate_burst":null,"failover":null}]

}
