package host

import (
	"fmt"
	"net/http"
)

/*
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r2 := new(Request)
	*r2 = *r
	r2.ctx = ctx
	return r2
}

*/

func ExampleTimeoutHandler() {

}

func _ExampleMiddleware() {
	m := http.NewServeMux()
	if m != nil {
	}
	//m.Handler()
	fmt.Printf("test () -> [%v]\n", "results")

	//Output:
	//fail
}
