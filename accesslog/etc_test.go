package accesslog

import "fmt"

func ExampleCreateHeaderEntry() {
	entry := createHeaderEntry(Reference{Operator: "", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "test", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "%REQ(", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "%REQ(t", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "%REQ()", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "%REQ(member)", Name: ""})
	fmt.Printf("Entry : %v\n", entry)
	entry = createHeaderEntry(Reference{Operator: "%REQ(member)", Name: "alias-member"})
	fmt.Printf("Entry : %v\n", entry)

	//Output:
	//Entry : {{ }  false}
	//Entry : {{ }  false}
	//Entry : {{ }  false}
	//Entry : {{ }  false}
	//Entry : {{ }  false}
	//Entry : {{header:member member}  true}
	//Entry : {{header:member alias-member}  true}
}

func _ExampleCreateEntry() {
	entry, err := createEntry(Reference{})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "static"})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)"})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)", Name: "new-name"})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC__%", Name: ""})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: ""})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: "new-name"})
	fmt.Printf("entry  : %v  Error  : %v\n", entry, err)

	//Output:
	//entry  : {   false}  Error  : invalid entry : operator is empty
	//entry  : {direct static  true}  Error  : <nil>
	//entry  : {header:static static  true}  Error  : <nil>
	//entry  : {header:static new-name  true}  Error  : <nil>
	//entry  : {   false}  Error  : invalid operator : operator not found or not a valid reference %TRAFFIC__%
	//entry  : {%TRAFFIC% traffic  true}  Error  : <nil>
	//entry  : {%TRAFFIC% new-name  true}  Error  : <nil>
}

func _ExampleCreateEntries() {
	var items []Entry

	err := CreateEntries(&items, []Reference{{Operator: "", Name: "name"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%INVALID", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	//err = addAttributes(&attrs, []Entry{{Operator: "static", Name: "name"}})
	//fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	//err = addAttributes(&attrs, []Entry{{Operator: "%START_TIME%", Name: ""}})
	//fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	//err = addAttributes(&attrs, []Entry{{Operator: "%START_TIME%", Name: "timestamp"}})
	//fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	//Output:
	//Attrs  : []  Error  : invalid entry : operator is empty
	//Attrs  : []  Error  : invalid operator : operator not found or not a valid reference %INVALID
}

//Attrs  : [{direct static name true}]  Error  : <nil>
//Attrs  : [{direct static name true} {%START_TIME% start_time  true}]  Error  : <nil>
//Attrs  : [{direct static name true} {%START_TIME% start_time  true} {%START_TIME% timestamp  true}]  Error

func ExampleCreateEntriesReq() {
	var items []Entry

	err := CreateEntries(&items, []Reference{{Operator: "%REQ(", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ()", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ(t", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ(customer)", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	//Output:
	//Entries  : []  Error  : invalid reference : operator is invalid %REQ(
	//Entries  : []  Error  : invalid reference : operator is invalid %REQ()
	//Entries  : []  Error  : invalid reference : operator is invalid %REQ(t
	//Entries  : [{{header:customer customer}  true}]  Error  : <nil>
}
