package accesslog

import (
	"fmt"
	"net/http"
)

func ExampleRequestHeaderTest() {
	req, _ := http.NewRequest("", "www.google.com", nil)
	req.Header.Add("customer", "Ted's Bait & Tackle")
	data := Logd{Req: req}
	fmt.Printf("Header : [%v] [%v]\n", "customer", data.Value(NewEntry("header:customer", "", "", true)))

	//Output:
	//Header : [customer] [Ted's Bait & Tackle]
}
