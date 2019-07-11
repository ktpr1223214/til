package main

import "fmt"

// cf. https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it
// go の map は参照ではない(というか、go に参照は存在しない)
// 何かと言うと、pointer to a runtime.hmap structure
func fn(m map[int]int) {
	fmt.Println(m == nil)
	m = make(map[int]int)
	fmt.Println(m == nil)
}

func modify(m map[int]int) {
	m[2] = 2
}

func main() {
	var m map[int]int
	fn(m)
	fmt.Println(m == nil)

	m = make(map[int]int)
	modify(m)
	fmt.Println(m[2])
}
