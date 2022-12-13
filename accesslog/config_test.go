package accesslog

import "fmt"

func translateEntry(e Entry) Entry {
	newE := Entry{Operator: e.Operator, Value: e.Value, Name: e.Name, StringValue: e.StringValue}
	if newE.Operator == "" {
		newE.Operator = "<empty>"
	}
	if newE.Value == "" {
		newE.Value = "<empty>"
	}
	if newE.Name == "" {
		newE.Name = "<empty>"
	}
	return newE
}

func Example_createHeaderEntry() {
	entry := createHeaderEntry(Reference{Operator: "", Name: ""})
	fmt.Printf("test: Entry(\"\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "test", Name: ""})
	fmt.Printf("test: Entry(\"test\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "%REQ(", Name: ""})
	fmt.Printf("test: Entry(\"REQ(\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "%REQ(t", Name: ""})
	fmt.Printf("test: Entry(\"REQ(t\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "%REQ()", Name: ""})
	fmt.Printf("test: Entry(\"REQ()\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "%REQ(member)", Name: ""})
	fmt.Printf("test: Entry(\"REQ(member)\") -> [%v]\n", translateEntry(entry))

	entry = createHeaderEntry(Reference{Operator: "%REQ(member)", Name: "alias-member"})
	fmt.Printf("test: Entry(\"REQ(member)\") -> [%v]\n", translateEntry(entry))

	//Output:
	//test: Entry("") -> [{<empty> <empty> <empty> false}]
	//test: Entry("test") -> [{<empty> <empty> <empty> false}]
	//test: Entry("REQ(") -> [{<empty> <empty> <empty> false}]
	//test: Entry("REQ(t") -> [{<empty> <empty> <empty> false}]
	//test: Entry("REQ()") -> [{<empty> <empty> <empty> false}]
	//test: Entry("REQ(member)") -> [{header:member member <empty> true}]
	//test: Entry("REQ(member)") -> [{header:member alias-member <empty> true}]

}

func Example_createEntry() {
	entry, err := createEntry(Reference{})
	fmt.Printf("test: createEntry({}) -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "static", Name: " "})
	fmt.Printf("test: createEntry(\"static\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC__%", Name: ""})
	fmt.Printf("test: createEntry(\"TRAFFIC__\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)"})
	fmt.Printf("test: createEntry(\"REQ(static)\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "%REQ(static)", Name: "new-name"})
	fmt.Printf("test: createEntry(\"REQ(static)\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: ""})
	fmt.Printf("test: createEntry(\"TRAFFIC\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	entry, err = createEntry(Reference{Operator: "%TRAFFIC%", Name: "new-name"})
	fmt.Printf("test: createEntry(\"TRAFFIC\") -> [%v] [err:%v]\n", translateEntry(entry), err)

	//Output:
	//test: createEntry({}) -> [{<empty> <empty> <empty> false}] [err:invalid entry reference : operator is empty ]
	//test: createEntry("static") -> [{<empty> <empty> <empty> false}] [err:invalid entry reference : name is empty [operator=static]]
	//test: createEntry("TRAFFIC__") -> [{<empty> <empty> <empty> false}] [err:invalid entry reference : operator not found %TRAFFIC__%]
	//test: createEntry("REQ(static)") -> [{header:static static <empty> true}] [err:<nil>]
	//test: createEntry("REQ(static)") -> [{header:static new-name <empty> true}] [err:<nil>]
	//test: createEntry("TRAFFIC") -> [{%TRAFFIC% traffic <empty> true}] [err:<nil>]
	//test: createEntry("TRAFFIC") -> [{%TRAFFIC% new-name <empty> true}] [err:<nil>]
}

func Example_CreateEntries() {
	var items []Entry

	err := CreateEntries(nil, []Reference{{Operator: "", Name: "name"}})
	fmt.Printf("test: CreateEntries(\"items: nil\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{})
	fmt.Printf("test: CreateEntries({}) -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "", Name: "name"}})
	fmt.Printf("test: CreateEntries(\"Operator: \"\"\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%INVALID", Name: ""}})
	fmt.Printf("test: CreateEntries(\"Operator: INVALID\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "static", Name: "name"}})
	fmt.Printf("test: CreateEntries(\"Operator: static\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%START_TIME%", Name: ""}})
	fmt.Printf("test: CreateEntries(\"Operator: START_TIME\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("test: CreateEntries(\"Operator: START_TIME\") -> [err:%v] [%v]\n", err, items)

	var newItems []Entry

	err = CreateEntries(&newItems, []Reference{{Operator: "%START_TIME%", Name: "timestamp"}, {Operator: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("test: CreateEntries(\"Operator: START_TIME\") -> [err:%v] [%v]\n", err, items)

	//Output:
	//test: CreateEntries("items: nil") -> [err:invalid configuration : entries are nil] [[]]
	//test: CreateEntries({}) -> [err:invalid configuration : configuration is empty] [[]]
	//test: CreateEntries("Operator: """) -> [err:invalid entry reference : operator is empty ] [[]]
	//test: CreateEntries("Operator: INVALID") -> [err:invalid entry reference : operator not found %INVALID] [[]]
	//test: CreateEntries("Operator: static") -> [err:<nil>] [[{direct static name true}]]
	//test: CreateEntries("Operator: START_TIME") -> [err:<nil>] [[{direct static name true} {%START_TIME% start_time  true}]]
	//test: CreateEntries("Operator: START_TIME") -> [err:<nil>] [[{direct static name true} {%START_TIME% start_time  true} {%START_TIME% timestamp  true}]]
	//test: CreateEntries("Operator: START_TIME") -> [err:invalid reference : name is a duplicate [timestamp]] [[{direct static name true} {%START_TIME% start_time  true} {%START_TIME% timestamp  true}]]

}

func _Example_CreateEntries_Request() {
	var items []Entry

	err := CreateEntries(&items, []Reference{{Operator: "%REQ(", Name: ""}})
	fmt.Printf("test: CreateEntries(\"REQ(\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ()", Name: ""}})
	fmt.Printf("test: CreateEntries(\"REQ()\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ(t", Name: ""}})
	fmt.Printf("test: CreateEntries(\"REQ(t)\") -> [err:%v] [%v]\n", err, items)

	err = CreateEntries(&items, []Reference{{Operator: "%REQ(customer)", Name: ""}})
	fmt.Printf("test: CreateEntries(\"REQ(customer)\") -> [err:%v] [%v]\n", err, items)

	//Output:
	//test: CreateEntries("REQ()") -> [err:invalid reference : operator is invalid %REQ(] [[]]
	//test: CreateEntries("REQ()") -> [err:invalid reference : operator is invalid %REQ()] [[]]
	//test: CreateEntries("REQ(t)") -> [err:invalid reference : operator is invalid %REQ(t] [[]]
	//test: CreateEntries("REQ(customer)") -> [err:<nil>] [[{header:customer customer  true}]]

}
