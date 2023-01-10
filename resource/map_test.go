package resource

import "fmt"

/*
func createTestEntry(uri string, status int32) *entry {
	entry := createEntry(uri, nil)
	entry.msgs.add(CreateMessage(VirtualHost, VirtualHost, StartupEvent, status, nil))
	return entry
}


*/
func ExampleDirectory_Add() {
	uri := "urn:test"
	uri2 := "urn:test:two"

	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 := directory.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directory.add(uri, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri)
	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 = directory.get(uri)
	fmt.Printf("test: get(%v) -> : %v\n", uri, d2)

	directory.add(uri2, nil)
	fmt.Printf("test: add(%v) -> : ok\n", uri2)
	fmt.Printf("test: count() -> : %v\n", directory.count())
	d2 = directory.get(uri2)
	fmt.Printf("test: get(%v) -> : %v\n", uri2, d2)

	fmt.Printf("test: uri() -> : %v\n", directory.uri())

	//Output:
	//test: count() -> : 0
	//test: get(urn:test) -> : <nil>
	//test: add(urn:test) -> : ok
	//test: count() -> : 1
	//test: get(urn:test) -> : &{urn:test <nil>}
	//test: add(urn:test:two) -> : ok
	//test: count() -> : 2
	//test: get(urn:test:two) -> : &{urn:test:two <nil>}
	//test: uri() -> : [urn:test urn:test:two]
	
}
