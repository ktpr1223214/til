package main

type Addifier interface{ Add(a, b int32) int32 }

type Adder struct{ name string }

//go:noinline
func (adder Adder) Add(a, b int32) int32 { return a + b }

// GOOS=linux GOARCH=amd64 go tool compile -m main.go
func main() {
	adder := Adder{name: "myAdder"}
	adder.Add(10, 32)           // doesn't escape
	Addifier(adder).Add(10, 32) // escapes なぜなら、interface が内部的に持つのはポインタのみなので、対象を heap で確保する必要があるから
}
