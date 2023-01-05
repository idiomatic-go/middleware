package actuator

import (
	"fmt"
	"net/http"
)

func ExampleMux_Error() {
	name := "google:get"
	m := newMux()

	err := m.add("", name)
	fmt.Printf("test: add() [err:%v] [count:%v]\n", err, m.count())

	err = m.add("https://google.com", "")
	fmt.Printf("test: add() [err:%v] [count:%v]\n", err, m.count())

	//Output:
	//test: add() [err:invalid configuration: pattern or name is empty [pattern:] [name:google:get]] [count:0]
	//test: add() [err:invalid configuration: pattern or name is empty [pattern:https://google.com] [name:]] [count:0]

}

func ExampleMux_Lookup() {
	google := "google:get"
	googleUri := "https://google.com"
	facebook := "facebook:get"
	facebookUri := "https://facebook.com"
	instagram := "instagram:get"
	instagramUri := "/instagram/pics"

	m := newMux()

	err := m.add(googleUri, google)
	err = m.add(facebookUri, facebook)
	err = m.add(instagramUri, instagram)
	fmt.Printf("test: add() [err:%v] [count:%v]\n", err, m.count())

	req, _ := http.NewRequest("", "https://twitter.com", nil)
	name, ok := m.lookup(req)
	fmt.Printf("test: lookup(%v) [name:%v] [ok:%v]\n", "https://twitter.com", name, ok)

	req, _ = http.NewRequest("", googleUri, nil)
	name, ok = m.lookup(req)
	fmt.Printf("test: lookup(%v) [name:%v] [ok:%v]\n", googleUri, name, ok)

	req, _ = http.NewRequest("", facebookUri, nil)
	name, ok = m.lookup(req)
	fmt.Printf("test: lookup(%v) [name:%v] [ok:%v]\n", facebookUri, name, ok)

	req, _ = http.NewRequest("", instagramUri, nil)
	name, ok = m.lookup(req)
	fmt.Printf("test: lookup(%v) [name:%v] [ok:%v]\n", instagramUri, name, ok)

	//Output:
	//test: add() [err:<nil>] [count:3]
	//test: lookup(https://twitter.com) [name:] [ok:false]
	//test: lookup(https://google.com) [name:google:get] [ok:true]
	//test: lookup(https://facebook.com) [name:facebook:get] [ok:true]
	//test: lookup(/instagram/pics) [name:instagram:get] [ok:true]

}
