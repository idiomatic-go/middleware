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

/*
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


*/
func Example_createOperator() {
	op, err := createOperator(accessdata.Operator{})
	fmt.Printf("test: createOperator({}) -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: " ", Value: "static"})
	fmt.Printf("test: createOperator(\"static\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "static", Value: "value"})
	fmt.Printf("test: createOperator(\"static\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%TRAFFIC__%"})
	fmt.Printf("test: createOperator(\"TRAFFIC__\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%REQ("})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%REQ(test"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%REQ()"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%REQ(static)"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "new-name", Value: "%REQ(static)"})
	fmt.Printf("test: createOperator(\"REQ(static)\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "", Value: "%TRAFFIC%"})
	fmt.Printf("test: createOperator(\"TRAFFIC\") -> [%v] [err:%v]\n", translateOperator(op), err)

	op, err = createOperator(accessdata.Operator{Name: "new-name", Value: "%TRAFFIC%"})
	fmt.Printf("test: createOperator(\"TRAFFIC\") -> [%v] [err:%v]\n", translateOperator(op), err)

	//Output:
	//test: createOperator({}) -> [{<empty> <empty>}] [err:invalid operator: value is empty ]
	//test: createOperator("static") -> [{<empty> <empty>}] [err:invalid operator: name is empty [static]]
	//test: createOperator("static") -> [{static value}] [err:<nil>]
	//test: createOperator("TRAFFIC__") -> [{<empty> <empty>}] [err:invalid operator: value not found or invalid %TRAFFIC__%]
	//test: createOperator("REQ(static)") -> [{<empty> <empty>}] [err:invalid operator: value not found or invalid %REQ(]
	//test: createOperator("REQ(static)") -> [{<empty> <empty>}] [err:invalid operator: value not found or invalid %REQ(test]
	//test: createOperator("REQ(static)") -> [{<empty> <empty>}] [err:invalid operator: value not found or invalid %REQ()]
	//test: createOperator("REQ(static)") -> [{static %REQ(static)}] [err:<nil>]
	//test: createOperator("REQ(static)") -> [{new-name %REQ(static)}] [err:<nil>]
	//test: createOperator("TRAFFIC") -> [{traffic %TRAFFIC%}] [err:<nil>]
	//test: createOperator("TRAFFIC") -> [{new-name %TRAFFIC%}] [err:<nil>]

}

func Example_CreateEntries() {
	var items []accessdata.Operator

	err := CreateOperators(nil, []accessdata.Operator{{Name: "name", Value: ""}})
	fmt.Printf("test: CreateOperators(\"items: nil\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{})
	fmt.Printf("test: CreateOperators({}) -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Name: "", Value: "%INVALID"}})
	fmt.Printf("test: CreateOperators(\"Value: INVALID\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Name: "name", Value: "static"}})
	fmt.Printf("test: CreateOperators(\"Value: static\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Name: "", Value: "%START_TIME%"}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	err = CreateOperators(&items, []accessdata.Operator{{Name: "duration", Value: "%DURATION%"}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	var newItems []accessdata.Operator

	err = CreateOperators(&newItems, []accessdata.Operator{{Name: "duration", Value: "%DURATION%"}, {Name: "duration", Value: "%DURATION%"}})
	fmt.Printf("test: CreateOperators(\"Value: START_TIME\") -> [err:%v] [%v]\n", err, items)

	//Output:
	//test: CreateOperators("items: nil") -> [err:invalid configuration: operators slice is nil] [[]]
	//test: CreateOperators({}) -> [err:invalid configuration: configuration slice is empty] [[]]
	//test: CreateOperators("Value: INVALID") -> [err:invalid operator: value not found or invalid %INVALID] [[]]
	//test: CreateOperators("Value: static") -> [err:<nil>] [[{name static}]]
	//test: CreateOperators("Value: START_TIME") -> [err:<nil>] [[{name static} {start_time %START_TIME%}]]
	//test: CreateOperators("Value: START_TIME") -> [err:<nil>] [[{name static} {start_time %START_TIME%} {duration %DURATION%}]]
	//test: CreateOperators("Value: START_TIME") -> [err:invalid operator: name is a duplicate [duration]] [[{name static} {start_time %START_TIME%} {duration %DURATION%}]]

}



