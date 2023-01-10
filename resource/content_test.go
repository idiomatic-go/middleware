package resource

import "fmt"

func ExampleAccessCredentials() {

	msg := Message{To: "to-uri", From: "from-uri", Content: []any{
		"text content",
		500,
		Credentials(func() (username, password string, err error) { return "", "", nil }),
	}}

	fmt.Printf("test: AccessCredentials(nil) -> %v\n", AccessCredentials(nil) != nil)
	fmt.Printf("test: AccessCredentials(mdg) -> %v\n", AccessCredentials(&Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", AccessCredentials(&msg) != nil)

	//Output:
	//test: AccessCredentials(nil) -> false
	//test: AccessCredentials(mdg) -> false
	//test: AccessCredentials(msg) -> true
}
