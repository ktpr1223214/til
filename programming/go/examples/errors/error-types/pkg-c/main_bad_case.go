// エラーの依存が悪い例
// 依存関係は、pkg-c <- pkg-b <- pkg-a だが、pkg-c でエラー処理の分岐をしたい場合、pkg-a を import する必要がある
package main

import (
	"fmt"
	"github.com/ktpr1223214/til/programming/go/errors/error-types/pkg-a"
	"github.com/ktpr1223214/til/programming/go/errors/error-types/pkg-b"
)

func main() {
	err := pkg_b.Fuga()
	switch err := err.(type) {
	case *pkg_a.MyError2:
		fmt.Println("MyError2", err)
	default:
		fmt.Println("Default")
	}
}
