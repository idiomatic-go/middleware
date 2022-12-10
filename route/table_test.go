package route

import "fmt"

func ExampleTable_SetDefault() {
	ti := NewTable()
	tt := ti.t()
	fmt.Printf("IsEmpty   : %v\n", tt.isEmpty())

	route := ti.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//ti.SetDefault(NewRoute("not-default"))
	//route = ti.Lookup(nil)
	//fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//Output:
	//IsEmpty   : true
	//IsDefault : true
	//IsDefault : false
}

func ExampleTable_Lookup() {
	ti := NewTable()
	tt := ti.t()
	fmt.Printf("IsEmpty   : %v\n", tt.isEmpty())

	route := ti.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//ti.SetDefault(NewRoute("not-default"))
	//route = ti.Lookup(nil)
	//fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//Output:
	//IsEmpty   : true
	//IsDefault : true
	//IsDefault : false
}
