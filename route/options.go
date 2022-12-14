package route

import "net/http"

// Matcher - provides matching a request to a route name in a Route table
type Matcher func(req *http.Request) (routeName string)
