package internal

import (
	"fmt"
	"io"
	"strings"
)

func ExampleStringReader() {
	s := "This is an example of content"
	r := ReaderCloser{strings.NewReader(s), nil}
	var buf = make([]byte, len(s))
	cnt, err := r.Read(buf)
	if cnt <= 0 || err != nil {
		fmt.Println("failure")
	} else {
		fmt.Println(fmt.Sprintf("Content: " + string(buf)))
	}

	//Output:
	//Content: This is an example of content
}

func ExampleFailure() {
	s := "This is an example of content"
	r := ReaderCloser{strings.NewReader(s), io.ErrUnexpectedEOF}
	var buf = make([]byte, len(s))
	_, err := r.Read(buf)
	if err == nil {
		fmt.Println("failure")
	} else {
		fmt.Printf("I/O Error: %v\n", err)
	}

	//Output:
	//I/O Error: unexpected EOF
}
