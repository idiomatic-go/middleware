package route

import "fmt"

func ExampleTable_SetDefault() {
	t := newTable()

	fmt.Printf("IsEmpty   : %v\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	t.SetDefault(NewRoute("not-default"))
	r = t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	//Output:
	//IsEmpty   : true
	//IsDefault : true
	//IsDefault : false
}

func ExampleTable_Lookup() {
	t := newTable()
	fmt.Printf("IsEmpty   : %v\n", t.isEmpty())

	r := t.Lookup(nil)
	fmt.Printf("IsDefault : %v\n", r.(*route).name == DefaultName)

	//ti.SetDefault(NewRoute("not-default"))
	//route = ti.Lookup(nil)
	//fmt.Printf("IsDefault : %v\n", route.IsDefault())

	//Output:
	//IsEmpty   : true
	//IsDefault : true

}
