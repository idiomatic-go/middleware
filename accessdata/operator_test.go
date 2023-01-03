package accessdata

import "fmt"

func _ExampleOperators() {
	op := Operators[DurationOperator]
	fmt.Printf("test: Operator() -> %v\n", op)

	op = Operators[StartTimeOperator]
	fmt.Printf("test: Operator() -> %v\n", op)

	//Output:
	//fail
}

func Example_IsRequestOperator() {
	op := Operator{}
	ok := IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(<empty>) -> %v\n", ok)

	op = Operator{Name: " ", Value: " "}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(<empty>) -> %v\n", ok)

	op = Operator{Name: "", Value: "REQ "}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(header"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ()"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(1)"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	op = Operator{Name: "", Value: "%REQ(header-name)"}
	ok = IsRequestOperator(op)
	fmt.Printf("test: IsRequestOperator(%v) -> %v\n", op, ok)

	//Output:
	//test: IsRequestOperator(<empty>) -> false
	//test: IsRequestOperator(<empty>) -> false
	//test: IsRequestOperator({ REQ }) -> false
	//test: IsRequestOperator({ %REQ(header}) -> false
	//test: IsRequestOperator({ %REQ()}) -> false
	//test: IsRequestOperator({ %REQ(1)}) -> true
	//test: IsRequestOperator({ %REQ(header-name)}) -> true

}

func Example_ParseRequestOperator() {
	op := Operator{}
	op2, ok := ParseRequestOperator(op)
	fmt.Printf("test: ParseRequestOperator() -> %v [op:%v] [op2:%v]\n", ok, op, op2)

	op = Operator{Name: "", Value: "%REQ("}
	op2, ok = ParseRequestOperator(op)
	fmt.Printf("test: ParseRequestOperator() -> %v [op:%v] [op2:%v]\n", ok, op, op2)

	op = Operator{Name: "", Value: "%REQ()"}
	op2, ok = ParseRequestOperator(op)
	fmt.Printf("test: ParseRequestOperator() -> %v [op:%v] [op2:%v]\n", ok, op, op2)

	op = Operator{Name: "", Value: "%REQ(1)"}
	op2, ok = ParseRequestOperator(op)
	fmt.Printf("test: ParseRequestOperator() -> %v [op:%v] [op2:%v]\n", ok, op, op2)

	op = Operator{Name: "", Value: "%REQ(name)"}
	op2, ok = ParseRequestOperator(op)
	fmt.Printf("test: ParseRequestOperator() -> %v [op:%v] [op2:%v]\n", ok, op, op2)

	//Output:
	//test: ParseRequestOperator() -> false [op:{ }] [op2:{ }]
	//test: ParseRequestOperator() -> false [op:{ %REQ(}] [op2:{ }]
	//test: ParseRequestOperator() -> false [op:{ %REQ()}] [op2:{ }]
	//test: ParseRequestOperator() -> true [op:{ %REQ(1)}] [op2:{1 %REQ(1)}]
	//test: ParseRequestOperator() -> true [op:{ %REQ(name)}] [op2:{name %REQ(name)}]

}
