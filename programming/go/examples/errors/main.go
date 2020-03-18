package main

import (
	"fmt"
	"log"
	
	"github.com/pkg/errors"
)

func Hoge() error {
	// errors.Errorf("hoge error") であれば、その時点での stack trace も記録される
	// cf. https://godoc.org/github.com/pkg/errors#Errorf
	// そのため、最下層が自身の定義する関数から errors.Errorf(...)を返す場合には、基本的に wrap は不要
	return fmt.Errorf("hoge error")
}

func Fuga() error {
	if err := Hoge(); err != nil {
		return errors.Wrap(err, "failed to execute Fuga")
	}
	return nil
}

func Piyo() error {
	if err := Fuga(); err != nil {
		// ここで、errors.Wrap(err, "failed to execute Piyo") とも出来るが
		// stack trace でここからの情報も別途表示される感じになるだけなので、stack trace 的な意味は（最下層でちゃんとしてれば）無いはず
		// ただし、文脈情報として、エラー message は追記されるため必要なこともあるかもしれない？？
		return errors.Wrap(err, "failed to execute Piyo")
	}
	return nil
}

func main() {
	log.Printf("Failed to Piyo: %s", Piyo())
	log.Printf("Failed to Piyo: %+v", Piyo())
}
