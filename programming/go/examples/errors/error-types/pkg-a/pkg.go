package pkg_a

import "errors"

var MyError = errors.New("my Error")

type MyError2 struct {
	Message string
}

func (m *MyError2) Error() string {
	return m.Message
}

func (m *MyError2) Temporary() bool {
	return true
}
