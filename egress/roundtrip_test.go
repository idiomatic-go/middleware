package egress

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"github.com/idiomatic-go/middleware/route"
	"net/http"
)

func Example_Routes() {
	r := Routes.Lookup(nil)
	fmt.Printf("test: Lookup(nil) -> [name:%v]\n", r.Name())

	//Output:
	//test: Lookup(nil) -> [name:*]
}

func Example_RoundTrip() {
	accesslog.SetEgressWrite(func(s string) {
		fmt.Printf("test: WriteEgress() -> [%v]\n", s)
	})
	name := "egress-route"
	r1, _ := route.NewRouteWithConfig(name, 1000, 500, 100, true, false)
	Routes.SetMatcher(func(req *http.Request) string {
		return name
	})
	Routes.Add(r1)
	w := wrapper{}

	resp, err := w.RoundTrip(nil)
	fmt.Printf("test: RoundTrip(nil) -> [resp:%v] [err:%v]\n", resp, err)

	req, _ := http.NewRequest("POST", "http://search.google.com/results", nil)
	resp, err = w.RoundTrip(req)
	fmt.Printf("test: RoundTrip(req) -> [resp:%v] [err:%v]\n", resp, err)

	EnableDefaultHttpClient()
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("test: RoundTrip(req) -> [resp:%v] [err:%v]\n", resp, err)

	//Output:
	//test: RoundTrip(nil) -> [resp:<nil>] [err:invalid argument : http.Request is nil on RoundTrip call]
	//test: RoundTrip(req) -> [resp:<nil>] [err:invalid wrapper : http.RoundTripper is nil]
}
