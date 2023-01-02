package ingress

import (
	"fmt"
	"net/http"
)

func ExampleMiddleware() {
	m := http.NewServeMux()
	if m != nil {
	}
	//m.Handler()
	fmt.Printf("test () -> [%v]\n", "results")

	//Output:
	//fail
}
