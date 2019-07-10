package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// 文字列を io.Reader にしたい場合の処理時間計測
const jsonStr = `{"username":"fugahoge","text":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadfef"}`

var a io.Reader
var aa *strings.Reader

var c []byte

func BenchmarkBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// memory allocation が2回走るのは、[]byte -> string の変換と interface への代入で実態が何かしらのポインタ型だからか
		a = bytes.NewBuffer([]byte(jsonStr))
		// c = []byte(jsonStr)
	}
}

func BenchmarkStringsReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// こちらも同様で、1回 memory allocation が走るのは、*string.Reader を使うことになるからかと
		// a = strings.NewReader(jsonStr)
		aa = strings.NewReader(jsonStr)
	}
}
