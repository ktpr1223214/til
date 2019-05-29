package pkg_b

import "github.com/ktpr1223214/til/programming/go/errors/error-types/pkg-a"

func Fuga() error {
	return &pkg_a.MyError2{
		Message: "Fuga",
	}
}
