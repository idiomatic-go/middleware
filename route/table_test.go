package route

import "fmt"

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("IsEmpty   : %v\n", t.isEmpty())

	route := t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", route.IsDefault())

	t.SetDefault(NewRoute("not-default"))
	route = t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//Output:
	//IsEmpty   : true
	//IsDefault : true
	//IsDefault : false
}

func ExampleTable_Lookup() {
	t := newTable()
	fmt.Printf("IsEmpty   : %v\n", t.isEmpty())

	route := t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//ti.SetDefault(NewRoute("not-default"))
	//route = ti.Lookup(nil)
	//fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//Output:
	//IsEmpty   : true
	//IsDefault : true
	//IsDefault : false
}
