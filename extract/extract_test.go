package extract

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accesslog"
	"net/http"
)

func _ExampleInitilizeUrl() {
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

func ExampleInifiniteLoop() {
	url = "localhost:8080/accesslog"

	req, _ := http.NewRequest("post", "localhost:8080/accesslog", nil)
	data := accesslog.Logd{Req: req}
	loop := sameUrl(&data)
	fmt.Printf("Infinite  : %v\n", loop)

	req, _ = http.NewRequest("post", "localhost:8081/accesslog", nil)
	data = accesslog.Logd{Req: req}
	loop = sameUrl(&data)
	fmt.Printf("Infinite  : %v\n", loop)

	//Output:
	//Infinite  : true
	//Infinite  : false

}
