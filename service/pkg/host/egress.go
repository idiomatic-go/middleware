package host

import (
	"encoding/json"
	"fmt"
	"github.com/idiomatic-go/middleware/actuator"
	"github.com/idiomatic-go/middleware/host"
	"github.com/idiomatic-go/middleware/service/pkg/resource"
	"net/http"
	"time"
)

const (
	googleNameFmt  = "fs/google/routes_%v.json"
	twitterNameFmt = "fs/twitter/routes_%v.json"
)

func initEgress() error {
	lookup := make(map[string]string, 16)

	host.WrapDefaultTransport()
	actuator.EgressTable.SetDefaultActuator(actuator.DefaultActuatorName, nil, actuator.NewTimeoutConfig(time.Second*4, http.StatusGatewayTimeout))

	routes, err := initEgressRoutes(fmt.Sprintf(googleNameFmt, GetRuntime()))
	if err != nil {
		return err
	}
	addEgressRoutes(lookup, routes)
	routes, err = initEgressRoutes(fmt.Sprintf(twitterNameFmt, GetRuntime()))
	if err != nil {
		return err
	}
	addEgressRoutes(lookup, routes)
	actuator.EgressTable.SetMatcher(func(req *http.Request) string {
		if req == nil || req.URL == nil {
			return ""
		}
		return lookup[req.URL.Host]
	})
	return nil
}

func readEgressRoutes(name string) ([]actuator.Route, error) {
	var routes []actuator.Route

	buf, err := resource.ReadFile(name)
	if err != nil {
		return routes, err
	}
	err1 := json.Unmarshal(buf, &routes)
	return routes, err1
}

func initEgressRoutes(name string) ([]actuator.Route, error) {
	routes, err := readEgressRoutes(name)
	if err != nil {
		return routes, err
	}
	for _, r := range routes {
		actuator.EgressTable.Add(r.Name, nil, r.Timeout, r.RateLimiter, r.Retry, r.Failover)
	}
	return routes, err
}

func addEgressRoutes(m map[string]string, routes []actuator.Route) {
	for _, r := range routes {
		if _, ok := m[r.Host]; !ok {
			m[r.Host] = r.Name
		}
	}
}
