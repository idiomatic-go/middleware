package host

import "fmt"

func Example_readRoutes() {
	name := "fs/google/routes_dev.json"
	routes, err := readRoutes(name)
	fmt.Printf("test: readRoutes(%v) -> [err:%v] [routes:%v]\n", name, err, len(routes))

	//Output:
	//fail

}
