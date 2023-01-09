package template

import (
	"errors"
	"fmt"
	"strings"
)

func lookupVariable(name string) (string, error) {
	switch strings.ToLower(name) {
	case "env":
		return "[ENV]", nil
	case "next":
		return "[NEXT]", nil
	case "last":
		return "[LAST]", nil
	}
	return "", errors.New(fmt.Sprintf("invalid argument : template variable is invalid: %v", name))
}

func ExampleExpandTemplateInvalidLookup() {
	// Lookup function is nil
	s := "test"
	_, err := ExpandTemplate(s, nil)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v\n", err)

	// Lookup name not found
	s = "test{invalid}"
	_, err = ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v\n", err)

	//Output:
	//Path Input  : test
	//Path Output : invalid argument : VariableLookup() is nil
	//Path Input  : test{invalid}
	//Path Output : invalid argument : template variable is invalid: invalid

}

func ExampleExpandTemplateInvalidDelimiters() {
	var err error
	// Mismatched delimiters - too many end delimiters
	s := "resources/test-file-name{env}}and{next}{last}.txt"
	_, err = ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v\n", err)

	// Mismatched delimiters - too many begin delimiters, this is valid as the extra begin delimiters are skipped
	s = "resources/test-file-name{env}and{next}{{last}.txt"
	path, err0 := ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v %v\n", path, err0)

	// Mismatched delimiters - embedded begin delimiter
	s = "resources/test-file-name{env}and{next{}{last}.txt"
	path, err0 = ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v %v\n", path, err0)

	//Output:
	//Path Input  : resources/test-file-name{env}}and{next}{last}.txt
	//Path Output : invalid argument : token has multiple end delimiters: env}}and
	//Path Input  : resources/test-file-name{env}and{next}{{last}.txt
	//Path Output : resources/test-file-name[ENV]and[NEXT][LAST].txt <nil>
	//Path Input  : resources/test-file-name{env}and{next{}{last}.txt
	//Path Output :  invalid argument : template variable is invalid:
}

func ExampleExpandTemplateValid() {
	s := ""
	path, err := ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : <empty> %v\n", err == nil)
	//fmt.Printf("Path Output : %v : %v\n", path, err)

	s = "resources/test-file-name-and-ext.txt"
	path, err = ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v : %v\n", path, err)

	s = "resources/test-file-name{env}and{next}{last}.txt"
	path, err = ExpandTemplate(s, lookupVariable)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v : %v\n", path, err)

	//Output:
	// Path Input  : <empty> true
	// Path Input  : resources/test-file-name-and-ext.txt
	// Path Output : resources/test-file-name-and-ext.txt : <nil>
	// Path Input  : resources/test-file-name{env}and{next}{last}.txt
	// Path Output : resources/test-file-name[ENV]and[NEXT][LAST].txt : <nil>

}
