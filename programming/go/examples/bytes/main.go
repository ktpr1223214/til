package main

import (
	"strings"
)

var d *strings.Reader

func main() {
	jsonStr := `{"username":"fugahoge","text":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadfef"}`

	// fmt.Println([]byte(jsonStr))
	// a = bytes.NewBuffer([]byte(jsonStr))
	d = strings.NewReader(jsonStr)
}
