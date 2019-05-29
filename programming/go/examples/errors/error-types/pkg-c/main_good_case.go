// エラーの依存が良い例
// 依存関係は、pkg-c <- pkg-b <- pkg-a で、pkg-c でエラー処理の分岐をしたい場合に interface のおかげで、pkg-a を import する必要がなくなっている
package main

import (
	"fmt"
	"github.com/ktpr1223214/til/programming/go/errors/error-types/pkg-b"
)

type temporary interface {
	Temporary() bool
}

func main() {
	err := pkg_b.Fuga()
	_, ok := err.(temporary)
	if ok {
		fmt.Println("MyError2", err)
	}
}
