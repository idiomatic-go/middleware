package vhost

import (
	"fmt"
)

func ExampleCreateToSend() {
	uriNone := "package:none"
	uriOne := "package:one"

	registerResourceUnchecked(uriNone, nil)
	registerResourceUnchecked(uriOne, nil)

	m := createToSend(nil, nil)
	fmt.Printf("Test {no override messages} : %v\n", m)

	em := MessageMap{"package:one": {To: "package:one", From: "fromUri", Event: StartupEvent}}
	m = createToSend(em, nil)
	fmt.Printf("Test {one override messages} : %v\n", m)

	//Output:
	// Test {no override messages} : map[package:none:{package:none event:startup vhost 0 []} package:one:{package:one event:startup vhost 0 []}]
	// Test {one override messages} : map[package:none:{package:none event:startup vhost 0 []} package:one:{package:one event:startup vhost 0 []}]
}

/*
func ExampleStatusUpdate() {
	uri := "progresql:main"

	registerPackageUnchecked(uri, nil)
	e := directory.get(uri)
	fmt.Printf("Entry : %v %v\n", e.uri, e.msgs.) //, e.statusChangeTS.Format(time.RFC3339))

	SendStartupSuccessfulResponse(uri)
	time.Sleep(time.Nanosecond * 1)
	e = directory.get(uri)
	fmt.Printf("Entry : %v %v\n", e.uri, e.startupStatus) //e.statusChangeTS.Format(time.RFC3339))

	//Output:
	// Entry : progresql:main 0
	// Entry : progresql:main 2

}

func ExampleValidateToSend() {
	uri := "package:none"

	registerPackageUnchecked(uri, nil)

	toSend := MessageMap{"invalid": {Event: StartupEvent, From: HostFrom}}
	err := validateToSend(toSend)
	fmt.Printf("Test - {invalid package uri in message} : %v\n", err)

	toSend = MessageMap{uri: {Event: StartupEvent, From: HostFrom}}
	err = validateToSend(toSend)
	fmt.Printf("Test - {valid package uri in message} : %v\n", err)

	uri2 := "package:one"
	registerPackageUnchecked(uri2, nil, []string{"package:invalid"})

	toSend = MessageMap{uri: {Event: StartupEvent, From: HostFrom}, uri2: {Event: StartupEvent, From: HostFrom}}
	err = validateToSend(toSend)
	fmt.Printf("Test - {invalid package uri in dependent} : %v\n", err)

	unregisterPackage(uri2)
	registerPackageUnchecked(uri2, nil, []string{"package:none"})

	toSend = MessageMap{"package:none": {Event: StartupEvent, From: HostFrom}, "package:one": {Event: StartupEvent, From: HostFrom}}
	err = validateToSend(toSend)
	fmt.Printf("Test - {valid package uri in dependent} : %v\n", err)

	//Output:
	// Test - {invalid package uri in message} : startup failure: directory entry does not exist for package uri: invalid
	// Test - {valid package uri in message} : <nil>
	// Test - {invalid package uri in dependent} : <nil>
	// Test - {valid package uri in dependent} : <nil>
}


*/
