package accessdata

import "fmt"

func ExampleOperators() {
	op := Operators[DurationOperator]
	fmt.Printf("test: Operator() -> %v\n", op)

	op = Operators[StartTimeOperator]
	fmt.Printf("test: Operator() -> %v\n", op)

	//Output:
	//fail
}
