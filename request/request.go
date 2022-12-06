package request

/*

func NewWithContext(ctx context.Context, url, method string, body io.Reader) (*http.Request, context.CancelFunc, error) {
	route := egress.LookupRoute(url, method)
	if route != nil && route.IsTimeout() {
		newCtx, cancel := context.WithTimeout(ctx, route.Timeout)
		req, err := http.NewRequestWithContext(newCtx, url, method, body)
		return req, cancel, err
	}
	req, err := http.NewRequestWithContext(ctx, url, method, body)
	return req, nil, err
}


*/
