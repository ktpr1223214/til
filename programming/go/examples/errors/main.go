package main

import (
	"github.com/pkg/errors"
	"log"
	"os"
)

// OpenFileNaive コンテキスト不明なエラーのサンプル実装
func OpenFileNaive(filePath string) error {
	_, err := os.Open(filePath)
	if err != nil {
		return err
	}
	return nil
}

// OpenFile コンテキスト情報付与バージョン
func OpenFile(filePath string) error {
	_, err := os.Open(filePath)
	if err != nil {
		// return errors.Wrapf(err, "os.Open failed at OpenFile")
		return errors.Wrapf(err, "failed to open file %s", filePath)
	}
	return nil
}

// SomeFunc OpenFile 呼び出しをする関数サンプル
// OpenFile 側で errors.Wrapf を呼び出しているため、ここで改めて Wrap しなくとも
// ファイル名・行数などの情報は stacktrace で得られる
func SomeFunc() error {
	if err := OpenFile("not_exist_sample.txt"); err != nil {
		return err
	}
	return nil
}

type SampleErr struct{}

func (s *SampleErr) Error() string {
	return "sample error"
}

func SampleError() error {
	var err *SampleErr
	return err
}

func SampleError2() error {
	var err error
	return err
}

func main() {
	// どっちかっていうと、interface の nil 話
	if err := SampleError(); err != nil {
		log.Println("error (nil of type *SampleErr)")
	}
	if err := SampleError2(); err != nil {
		log.Println("error (nil, nil)")
	}

	// コンテキスト(どの関数が os.Open で失敗したのか)がわからない
	if err := OpenFileNaive("not_exist_file.txt"); err != nil {
		log.Println(err)
	}

	// コンテキスト情報あり
	if err := OpenFile("not_exist_file.txt"); err != nil {
		log.Printf("%v", err)
	}

	// コンテキスト情報あり
	// %+v で stacktrace(ファイル名・行数)
	if err := SomeFunc(); err != nil {
		log.Printf("%+v", err)
	}
}
