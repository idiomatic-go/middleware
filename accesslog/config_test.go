package accesslog

import (
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
)

func translateOperator(op accessdata.Operator) accessdata.Operator {
	newOp := accessdata.Operator{Name: op.Name, Value: op.Value}
	if newOp.Name == "" {
		newOp.Name = "<empty>"
	}
	if newOp.Value == "" {
		newOp.Value = "<empty>"
	}
	//if newE.Name == "" {
	//	newE.Name = "<empty>"
	//}
	return newOp
}

func Example_createHeaderOperator() {
	op := createHeaderOperator(accessdata.Operator{Name: "", Value: ""})
	fmt.Printf("test: Operator(\"\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "test", Name: ""})
	fmt.Printf("test: Operator(\"test\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "%REQ(", Name: ""})
	fmt.Printf("test: Operator(\"REQ(\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "%REQ(t", Name: ""})
	fmt.Printf("test: Operator(\"REQ(t\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "%REQ()", Name: ""})
	fmt.Printf("test: Operator(\"REQ()\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "%REQ(member)", Name: ""})
	fmt.Printf("test: Operator(\"REQ(member)\") -> [%v]\n", translateOperator(op))

	op = createHeaderOperator(accessdata.Operator{Value: "%REQ(member)", Name: "alias-member"})
	fmt.Printf("test: Operator(\"REQ(member)\") -> [%v]\n", translateOperator(op))

	//Output:
	//test: Operator("") -> [{<empty> <empty>}]
	//test: Operator("test") -> [{<empty> <empty>}]
	//test: Operator("REQ(") -> [{<empty> <empty>}]
	//test: Operator("REQ(t") -> [{<empty> <empty>}]
	//test: Operator("REQ()") -> [{<empty> <empty>}]
	//test: Operator("REQ(member)") -> [{member header:member}]
	//test: Operator("REQ(member)") -> [{alias-member header:member}]

}

func Example_createOperator() {
	op, err := createOperator(accessdata.Operator{})
	fmt.Printf("test: createOperator({}) -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "static", Name: " "})
	fmt.Printf("test: createOperator(\"static\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "%TRAFFIC__%", Name: ""})
	fmt.Printf("test: createOperator(\"TRAFFIC__\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "%REQ(static)"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "%REQ(static)", Name: "new-name"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "%TRAFFIC%", Name: ""})
	fmt.Printf("test: createOperator(\"TRAFFIC\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Value: "%TRAFFIC%", Name: "new-name"})
	fmt.Printf("test: createOperator(\"TRAFFIC\") -> [%v] [err:%v]\n", translateOperator(op), err)

	//Output:
	//test: createOperator({}) -> [{<empty> <empty>}] [err:invalid operator: value is empty ]
	//test: createOperator("static") -> [{<empty> <empty>}] [err:invalid operator : name is empty [static]]
	//test: createOperator("TRAFFIC__") -> [{<empty> <empty>}] [err:invalid operator : value not found %TRAFFIC__%]
	//test: createOperator("REQ(static)") -> [{static header:static}] [err:<nil>]
	//test: createOperator("REQ(static)") -> [{new-name header:static}] [err:<nil>]
	//test: createOperator("TRAFFIC") -> [{traffic %TRAFFIC%}] [err:<nil>]
	//test: createOperator("TRAFFIC") -> [{new-name %TRAFFIC%}] [err:<nil>]

}

func Example_CreateEntries() {
	var items []accessdata.Operator

	err := CreateOperators(nil, []accessdata.Operator{{Value: "", Name: "name"}})
	fmt.Printf("test: CreateOperators(\"items: nil\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{})
	fmt.Printf("test: CreateOperators({}) -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Value: "", Name: "name"}})
	fmt.Printf("test: CreateOperators(\"Value: \"\"\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Value: "%INVALID", Name: ""}})
	fmt.Printf("test: CreateOperators(\"Value: INVALID\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Value: "static", Name: "name"}})
	fmt.Printf("test: CreateOperators(\"Value: static\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Value: "%START_TIME%", Name: ""}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Value: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	var newItems []accessdata.Operator

	err = CreateOperators(&newItems, []accessdata.Operator{{Value: "%START_TIME%", Name: "timestamp"}, {Value: "%START_TIME%", Name: "timestamp"}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	//Output:
	//test: CreateOperators("items: nil") -> [err:invalid configuration : operators are nil] [[]]
	//test: CreateOperators({}) -> [err:invalid configuration : configuration is empty] [[]]
	//test: CreateOperators("Value: """) -> [err:invalid operator: value is empty name] [[]]
	//test: CreateOperators("Value: INVALID") -> [err:invalid operator : value not found %INVALID] [[]]
	//test: CreateOperators("Value: static") -> [err:<nil>] [[{direct:name static}]]
	//test: CreateOperators("Value: START_TIME") -> [err:<nil>] [[{direct:name static} {start_time %START_TIME%}]]
	//test: CreateOperators("Value: START_TIME") -> [err:<nil>] [[{direct:name static} {start_time %START_TIME%} {timestamp %START_TIME%}]]
	//test: CreateOperators("Value: START_TIME") -> [err:invalid reference : name is a duplicate [timestamp]] [[{direct:name static} {start_time %START_TIME%} {timestamp %START_TIME%}]]

}

/*
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

*/
