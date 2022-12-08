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
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "static", Name: " "})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)"})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)", Name: "new-name"})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC__%", Name: ""})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: ""})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: "new-name"})
	fmt.Printf("Entry  : %v  Error  : %v\n", entry, err)

	//Output:
	//Entry  : {{ }  false}  Error  : invalid entry reference : operator is empty
	//Entry  : {{ }  false}  Error  : invalid entry reference : name is empty [operator=static]
}

//Entry  : {{ }  false}  Error  : invalid entry reference : operator is empty
//Entry  : {{ }  false}  Error  : invalid entry reference : name is empty [operator=static]
//Entry  : {{header:static static}  true}  Error  : <nil>
//Entry  : {{header:static new-name}  true}  Error  : <nil>
//Entry  : {{ }  false}  Error  : invalid entry reference : operator not found %TRAFFIC__%
//Entry  : {{%TRAFFIC% traffic}  true}  Error  : <nil>
//Entry  : {{%TRAFFIC% new-name}  true}  Error  : <nil>

func ExampleCreateEntries() {
	var items []Entry

	err := CreateEntries(nil, []Reference{{Operator: "", Name: "name"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "", Name: "name"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%INVALID", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "static", Name: "name"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%START_TIME%", Name: ""}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	err = CreateEntries(&items, []Reference{{Operator: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", items, err)

	var newItems []Entry

	err = CreateEntries(&newItems, []Reference{{Operator: "%START_TIME%", Name: "timestamp"}, {Operator: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("Entries  : %v  Error  : %v\n", newItems, err)

	//Output:
	//Entries  : []  Error  : invalid configuration : entry slice is nil
	//Entries  : []  Error  : invalid configuration : configuration is empty
	//Entries  : []  Error  : invalid entry reference : operator is empty
	//Entries  : []  Error  : invalid entry reference : operator not found %INVALID
	//Entries  : [{{direct static} name true}]  Error  : <nil>
	//Entries  : [{{direct static} name true} {{%START_TIME% start_time}  true}]  Error  : <nil>
	//Entries  : [{{direct static} name true} {{%START_TIME% start_time}  true} {{%START_TIME% timestamp}  true}]  Error  : <nil>
	//Entries  : [{{%START_TIME% timestamp}  true}]  Error  : invalid reference : name is a duplicate [timestamp]
}

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
