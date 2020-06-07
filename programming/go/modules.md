---
title: Modules
---

## 概要
* Semantic Versioning に基づいたバージョン管理
  * v3.1.4:
    * Major version: increment for backwards-incompatible changes
    * Minor version: increment for new features
    * Patch version: increment for bug fixes
* モジュールは、リポジトリのバージョンタグ/リビジョン毎で管理
* module-aware mode と GOPATH mode が存在
  * GOPATH mode では、```$GOPATH/src``` 以下からサードパーティのパッケージを読み込む
  * module-aware mode では、```$GOPATH/pkg/mod/``` 以下にバージョン毎のパッケージの実体がおかれる
  * GO111MODULE という環境変数で切り替わる
    * on: module-aware mode
    * off: GOPATH mode
    * auto: $GOPATH 配下では GOPATH mode / それ以外のディレクトリでは、module-aware mode
* go.mod
  * module: ルートディレクトリのモジュール名
  * require: 必要なモジュール名とバージョン名を指定
  * exclude: 明示的に除外するモジュールを指定
  * replace: require で指定したモジュール名を置き換える

### 基本
* package
  * Go programs are constructed by linking together packages. A package in turn is constructed from one or more source files that together declare constants, types, variables and functions belonging to the package and which are accessible in all files of the same package. Those elements may be exported and used in another package
* module
  * A module is a collection of Go packages stored in a file tree with a go.mod file at its root.
* go.mod ファイルの役割
  * The go.mod file defines the module’s module path, which is also the import path used for the root directory, and its dependency requirements, which are the other modules needed for a successful build. Each dependency requirement is written as a module path and a specific semantic version
  * Packages in subdirectories have import paths consisting of the module path plus the path to the subdirectory
  * Only direct dependencies are recorded in the go.mod file
    * 勿論、indirect な dependency もあり、```go list -m all``` で出力可
  * peuso-version でバージョン指定される場合もある
    * https://golang.org/cmd/go/#hdr-Pseudo_versions
  * The go.mod file is meant to be readable and editable by both programmers and tools
  * ```go list -m -versions <module>``` で利用可能な一覧を取得
  * ```go get <module>@v...``` でバージョン指定して、更新（デフォルトは、@latest）
  * あるモジュールの利用がコードからなくなっても、```go test``` なりを実行しても、その依存が削除されることはない
    * Because building a single package, like with go build or go test, can easily tell when something is missing and needs to be added, but not when something can safely be removed. Removing a dependency can only be done after checking all packages in a module, and all possible build tag combinations for those packages
    * これをやるには、```go mod tidy``` の実行が必要
* When it encounters an import of a package not provided by any module in go.mod, the go command automatically looks up the module containing that package and adds it to go.mod, using the latest version. (“Latest” is defined as the latest tagged stable (non-prerelease) version, or else the latest tagged prerelease version, or else the latest untagged version.)
* go.sum ファイルの役割
  * go.sum containing the expected cryptographic hashes of the content of specific module versions
  * The go command uses the go.sum file to ensure that future downloads of these modules retrieve the same bits as the first download, to ensure the modules your project depends on do not change unexpectedly, whether for malicious, accidental, or other reasons.
  * format
    * ```<module> <version>[/go.mod] <hash>```
    * モジュールのバージョン毎に 2 行生成される
      * The first line gives the hash of the module version's file tree
      * The second line appends "/go.mod" to the version and gives the hash of only the module version's (possibly synthesized) go.mod file
        * The go.mod-only hash allows downloading and authenticating a module version's go.mod file, which is needed to compute the dependency graph, without also downloading all the module's source code
    * h1: SHA-256

### command
``` bash
# go help list から抜粋
# The arguments to list -m are interpreted as a list of modules, not packages.
# The main module is the module containing the current directory.
# The active modules are the main module and its dependencies.
# With no arguments, list -m shows the main module.
# With arguments, list -m shows the modules specified by the arguments.
# Any of the active modules can be specified by its module path.
# The special pattern "all" specifies all the active modules, first the main
# module and then dependencies sorted by module path.

# main module の確認
$ go list -m

# <main module> は上の結果
# main module のパッケージ一覧
$ go list <main module>/...

# all the active modules
$ go list -m all

# include the transitive dependencies of the named packages
$ go list -deps <pkg path>

# main module のパッケージから、<module> で指定したモジュールのパッケージへの最短パスを見つける
# The main module is the module containing the current directory.(go list -m)
# <module> を指定して、返すのはパッケージであることに注意
$ go mod why -m <module>

# <package> 指定で、上と同じ条件でパッケージを返す
$ go mod why <package>
```

### go mod list
[Using go list, go mod why and go mod graph](https://github.com/go-modules-by-example/index/blob/master/018_go_list_mod_graph_why/README.md) の理解補足

* Firstly github.com/davecgh/go-spew: の部分
  * パッケージの transitive dependencies の中では、github.com/davecgh/go-spew/spew を import している部分がない、ってことなのか
  * 素直に見ると、github.com/sirupsen/logrus.test から assert の流れで依存なので、テスト含めた形でのモジュール依存はしているが、パッケージの transitive dependencies にそれが現れてこない、っていうのはそれはそうな気がする

* Secondly, github.com/kr/pty: の部分
  * This tells us that there is no package from the github.com/kr/pty module in the transitive import graph of of the main module
    * つまり、main module 内部の package import をみていった場合に、github.com/kr/pty module のパッケージは import していないということっぽい
    * The main module is the module containing the current directory.
  * go.mod ファイルからは
    * github.com/kr/pty <- github.com/kr/text <- github.com/kr/pretty がわかる
  * つまり、github.com/kr/text は package import から辿れてその依存から github.com/kr/pty も導かれるということ
    * package import としては辿れずとも、ということ

## Reference
* [Using Go Modules](https://blog.golang.org/using-go-modules)
* [Semantic Versioning 2.0.0](https://semver.org/#spec-item-9)
  * pre-release の付け方など
* [Using go list, go mod why and go mod graph](https://github.com/go-modules-by-example/index/blob/master/018_go_list_mod_graph_why/README.md)
  * 表題全般の詳細例
  * indirect についてもわかりやすい
* [Our Software Dependency Problem](https://research.swtch.com/deps)
  * In today’s software development world, a dependency is additional code that you want to call from your program. Adding a dependency avoids repeating work already done: designing, writing, testing, debugging, and maintaining a specific unit of code
  * A dependency manager (sometimes called a package manager) automates the downloading and installation of dependency packages
  * We are trusting more code with less justification for doing so
  * inspecting a package and deciding whether to depend on it の観点
    * Design
    * Code Quality
    * Testing
    * Debugging
    * Maintenance
    * Usage
    * Security
    * Licensing
    * Dependencies