package resource

import (
	"fmt"
)

func ExampleCreateToSend() {
	none := "/resource/none"
	one := "/resource/one"

	registerResourceUnchecked(none, nil)
	registerResourceUnchecked(one, nil)

	m := createToSend(nil, nil)
	fmt.Printf("test: registerResourceUncheckde(nil,nil) -> %v\n", m)

	cm := ContentMap{one: "test content"}
	m = createToSend(cm, nil)
	fmt.Printf("test: registerResourceUncheckde(map,nil) -> %v\n", m)

	//Output:
	//test: registerResourceUncheckde(nil,nil) -> map[/resource/none:{/resource/none host event:startup -101 <nil> <nil>} /resource/one:{/resource/one host event:startup -101 <nil> <nil>}]
	//test: registerResourceUncheckde(map,nil) -> map[/resource/none:{/resource/none host event:startup -101 <nil> <nil>} /resource/one:{/resource/one host event:startup -101 test content <nil>}]

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
