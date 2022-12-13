package accesslog

import (
	"fmt"
	"strings"
)

func Example_WriteJson() {
	sb := strings.Builder{}

	writeJson(&sb, "first", "string value", true)
	writeJson(&sb, "second", "100", false)
	writeJson(&sb, "third", "another string value", true)
	writeJson(&sb, "fourth", "true", false)
	sb.WriteString("}")

	fmt.Printf("test: writeJson() -> [%v]\n", sb.String())

	//Output:
	//test: writeJson() -> [{"first":"string value","second":100,"third":"another string value","fourth":true}]
}
