package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	name string
}

// for stringer interface
func (p *Person) String() string {
	return p.name
}

type T struct{}

func (t T) F() {}

type P interface {
	F()
}

func newT() *T { return new(T) }

type Thing struct {
	P
}

func factory(p P) *Thing {
	return &Thing{P: p}
}

const ENABLE_FEATURE = false

type doError struct{}

func (d *doError) Error() string {
	return "doError"
}

func do() error {
	var err *doError
	return err
}

func do2() *doError {
	return nil
}

func wrapDo() error {
	return do2()
}

func main() {
	fmt.Println("Interface init:")
	var fuga interface{}
	fmt.Println(reflect.TypeOf(fuga))         // <nil>
	fmt.Println(reflect.ValueOf(fuga).Kind()) // invalid
	fmt.Println("")

	fmt.Println("Slice:")
	var s []int
	fmt.Println(reflect.TypeOf(s))          // # []int
	fmt.Println(reflect.ValueOf(s))         // []
	fmt.Println(reflect.ValueOf(s).IsNil()) // true
	fmt.Println(reflect.ValueOf(s).Kind())  // slice
	fmt.Println("")

	fmt.Println("Interface and nil:")
	// nil は型を持たないので、
	// a := nil #  use of untyped nil と怒られる

	// これは pointer
	fmt.Println("Pointer:")
	var p *Person
	fmt.Println(reflect.TypeOf(p))  // *main.Person
	fmt.Println(reflect.ValueOf(p)) // nil
	fmt.Println(p == nil)
	fmt.Println("")

	fmt.Println("Interface:")
	var ss fmt.Stringer
	fmt.Println(reflect.TypeOf(ss))  // <nil>
	fmt.Println(reflect.ValueOf(ss)) // <invalid reflect.Value>
	ss = p
	// ここで interface ss は型を持つ
	fmt.Println(reflect.TypeOf(ss))          // *main.Person
	fmt.Println(reflect.ValueOf(ss))         // <nil>
	fmt.Println(reflect.ValueOf(ss).IsNil()) // true
	// ss の value は nil だが、type が *main.Person なので ss == nil は成立していない
	// interface は (type, value) で規定され、いずれもが(nil, nil)でなければ nil ではない
	fmt.Println(ss == nil)
	fmt.Println("")

	// cf. https://speakerdeck.com/campoy/understanding-nil?slide=55
	fmt.Println("Error1(Do not declare concrete error vars):")
	// 返ってくる err(interface)は、(*doError, nil)
	fmt.Println(do() == nil)          // false
	fmt.Println(reflect.TypeOf(do())) // *main.doError
	fmt.Println("")

	fmt.Println("Error2(do not return concrete error vars):")
	// 返ってくる err(interface)は、(*doError, nil)
	fmt.Println(do2() == nil)             // true
	fmt.Println(wrapDo() == nil)          // false
	fmt.Println(reflect.TypeOf(wrapDo())) // *main.doError
	fmt.Println("")

	// another example
	// cf. https://dave.cheney.net/2017/08/09/typed-nils-in-go-2
	fmt.Println("Anothe Example:")
	t := newT()
	t2 := t
	if !ENABLE_FEATURE {
		t2 = nil
	}
	thing := factory(t2)
	fmt.Println(thing.P == nil) // false
	// 解説:
	// nil is a compile time constant which is converted to whatever type is required,
	// just as constant literals like 200 are converted to the required integer type automatically.

	// Given the expression p == nil, both arguments must be of the same type,
	// therefore nil is converted to the same type as p, which is an interface type.
	// So we can rewrite the expression as (*T, nil) == (nil, nil).

	// As equality in Go almost always operates as a bitwise comparison it is clear that the memory bits which hold the interface value (*T, nil)
	// are different to the bits that hold (nil, nil) thus the expression evaluates to false.
}
