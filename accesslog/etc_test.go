package accesslog

import "fmt"

func ExampleParseHeaderAttribute() {
	attr := parseHeaderAttribute(Entry{Operator: "", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "test", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "%REQ(", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "%REQ(t", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "%REQ()", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "%REQ(member)", Name: ""})
	fmt.Printf("Attr : %v\n", attr)
	attr = parseHeaderAttribute(Entry{Operator: "%REQ(member)", Name: "alias-member"})
	fmt.Printf("Attr : %v\n", attr)

	//Output:
	//Attr : {   false}
	//Attr : {   false}
	//Attr : {   false}
	//Attr : {   false}
	//Attr : {   false}
	//Attr : {header:member member  true}
	//Attr : {header:member alias-member  true}
}

func _ExampleCreateAttribute() {
	attr, err := createAttribute(Entry{})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "static"})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "%REQ(static)"})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "%REQ(static)", Name: "new-name"})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "%TRAFFIC__%", Name: ""})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "%TRAFFIC%", Name: ""})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	attr, err = createAttribute(Entry{Operator: "%TRAFFIC%", Name: "new-name"})
	fmt.Printf("Attr  : %v  Error  : %v\n", attr, err)

	//Output:
	//Attr  : {   false}  Error  : invalid entry : operator is empty
	//Attr  : {direct static  true}  Error  : <nil>
	//Attr  : {header:static static  true}  Error  : <nil>
	//Attr  : {header:static new-name  true}  Error  : <nil>
	//Attr  : {   false}  Error  : invalid operator : operator not found or not a valid reference %TRAFFIC__%
	//Attr  : {%TRAFFIC% traffic  true}  Error  : <nil>
	//Attr  : {%TRAFFIC% new-name  true}  Error  : <nil>
}

func _ExampleAddAttributes() {
	var attrs []attribute

	err := addAttributes(&attrs, []Entry{{Operator: "", Name: "name"}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	err = addAttributes(&attrs, []Entry{{Operator: "%INVALID", Name: ""}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

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

func ExampleAddAttributesReq() {
	var attrs []attribute

	err := addAttributes(&attrs, []Entry{{Operator: "%REQ(", Name: ""}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	err = addAttributes(&attrs, []Entry{{Operator: "%REQ()", Name: ""}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	err = addAttributes(&attrs, []Entry{{Operator: "%REQ(t", Name: ""}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	err = addAttributes(&attrs, []Entry{{Operator: "%REQ(customer)", Name: ""}})
	fmt.Printf("Attrs  : %v  Error  : %v\n", attrs, err)

	//Output:
	//Attrs  : []  Error  : invalid entry : operator is invalid %REQ(
	//Attrs  : []  Error  : invalid entry : operator is invalid %REQ()
	//Attrs  : []  Error  : invalid entry : operator is invalid %REQ(t
	//Attrs  : [{header:customer customer  true}]  Error  : <nil>
}
