---
title: Error
---

## Error

### error 種別
* dave 参照

1. sentinel error
    * エラーを値としてもつ
``` go
// cf. io pkg
var EOF = errors.New("EOF")

// how to use
if err == ErrSomething { … }
```
* 等号比較が必要
    * fmt.Errorf などでコンテキスト付与をするとこれが崩れる
    * よって、Error() で文字列の中を見たりする必要が出てくるが、これは推奨されない
* Sentinel errors become part of your public API
    * エラー値を元にした分岐を別のパッケージで使い始めたりすると、更にそれに依存したコードが...となり、
    元のエラー側の変更コストが大きくなる
* Sentinel errors create a dependency between two packages
    * 1つのパッケージで考えるよりは、こういうパッケージが複数あるプロジェクトを想像して、
    その場合に利用側でこれらパッケージ郡を import しないといけなくなるというのが分かり易いか

2. error types
    * error interface を実装した構造体
``` go
type MyError struct {
        Msg string
        File string
        Line int
}

func (e *MyError) Error() string {
        return fmt.Sprintf("%s:%d: %s”, e.File, e.Line, e.Msg)
}

// how to use
err := something()
switch err := err.(type) {
case nil:
        // call succeeded, nothing to do
case *MyError:
        fmt.Println(“error occurred on line:”, err.Line)
default:
// unknown error
}
```

* コンテキスト付与が可能に
    * wrap an underlying error to provide more context
``` go
// cf. os.PathError
type PathError struct {
	Op   string
	Path string
	// the cause
	Err  error
}

// 使う時
return &PathError{"readat", f.name, errors.New("negative offset")}
```
* If your code implements an interface whose contract requires a specific error type,
all implementors of that interface need to depend on the package that defines the error type.

3. opaque errors
    * interface をパッケージ内に持つ
    * I call this style opaque error handling,
    because while you know an error occurred,
    you don’t have the ability to see inside the error.
    As the caller, all you know about the result of the operation is that it worked, or it didn’t.

* 単に、エラーの有無というのでは足りないことがある
    * ex. ネットワークエラーの種別に応じてリトライ可否を判断したい

``` go
type temporary interface {
        Temporary() bool
}

// IsTemporary returns true if err is temporary.
func IsTemporary(err error) bool {
        te, ok := err.(temporary)
        return ok && te.Temporary()
}
```

* cf. net.Error

``` go
type Error interface {
        error
        Timeout() bool   // Is the error a timeout?
        Temporary() bool // Is the error temporary?
}
```

* The key here is this logic can be implemented without importing the package that defines the error
or indeed knowing anything about err‘s underlying type–we’re simply interested in its behaviour.

#### もう少し考えると
* A <- B <- C(ここで error 定義)というパッケージの依存があるとする
    * B はかならず C を import していることになる
    * このときに、A がどうなるのか？を考える
* 1 or 2のエラーの場合
    * A で B が返すエラーの種別により処理を分岐したい場合、C を import する必要が出てくる？
* 3 の場合
    * A は C を必ずしも import せずとも、必要な interface 定義で対応可能？

### error 返り値について
``` go
type MyErr struct {}

func (m *MyErr) Error() string {
    ...
}

func Hoge() error {
    ...
    return &MyErr{}
}
```
* 暗黙的に代入が発生する場所がある
    * 関数呼び出しは、引数の値を対応するパラメータ変数に暗黙的に代入
    * return 文は return のオペランドを対応する結果変数へ暗黙的に代入
* error の場合は後者のパターン

## Reference
* [Don’t just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
    * 3種類のエラーの話
        * どっちかっていうと、pkg 向けの話なのでアプリだとまた別の観点で評価出来そうな事には留意
    * Don’t just check errors, handle them gracefully
    * Only handle errors once
* [Failure is your Domain](https://middlemost.com/failure-is-your-domain/)
    * github.com/pkg/errors とはまた違う観点から