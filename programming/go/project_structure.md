---
title: Project Structure
---

## Project Structure

### 何故考えるべきか
* Go をプロジェクトで使う場合に、中長期なメンテのコスト・変更の容易さを担保したい
* そのためにはある程度最初から考慮すべき
    * 後から構造を変えるのは結構大変     

### What is good structure?
* 一貫性がある
* 理解が容易
* 変更が容易
* メンテナンスが容易    
    * この辺りはプロジェクトの複雑性を下げるために必須
* テストが容易
* 構造に設計が反映されている
    * DDD 的な話で、設計とコードが対応していてそれが構造からもわかる的なイメッジ？

## アプローチ
### Flat Structure
``` 
├── main.go
├── handler.go
├── model.go
├── storage.go
└── ...
```
* pros
    * シンプル
        * circular dependency も起き得ない
    * 最初はここから初めて、package 切っていくというのはオススメ        
* cons        
    * 規模がある程度を超えると厳しい
        * 中身を一々見ないと何をしているのかわからない
        * それを嫌がって、1フォルダに数十とかファイルを作っていくか？ → 厳しい
    * グローバルな変数使われると(使うなって話ではあるが)全体に波及してしまう
        * 同じ package なので言語的に食い止めることはできない
        
* flat からまとめる部分をつくっていくのが以下

### Layered Architecture
* Group By Function
    * presentation / user interface
    * business logic
    * external dependencies / infrastructure
* pros
    * どこに何を置くかは明確になる
* cons
    * (go では)naming convention に合わない
    * circular dependency になりやすい

### Group By Module
* Group By Module
    * ex. beers / database / storage
        * 例えば handler なども、beers に関わるものはすべて同じ Module として扱う 
* cons
    * (go では)naming convention に合わない
        * ex. beers.Beer / storage.Storage など
    * 何をどの module に置くべきかがそこまで明確でない
        * ex. ビールをレビューするサービスで、beers/reviews は分けるべきなのか、それとも reviews は beers の submodule とすべきか
    * 機能追加や機能を使うのも難しい
        * 上と似ているが、ビールレビュー追加っていうのは、review 側なのか beer 側にあるのか             
    * circular dependency になりやすい

### Group By Context
* Group By Context
    * 要は DDD
    * package の分け方は、その package が providing するものに基づく(package に含まれるものではなく)   
        * ex. reviewing / listing /adding
        * これは何ソース？
            * https://www.citerus.se/go-ddd
            * この話かと
* pros
    * package 名をみると、何のアプリなのかよくわかる
    * 機能追加や利用も package 名からわかりやすい
        * 上の例でいうと、ビールレビュー追加は adding 配下になるはず
    * circular dependency になりにくい
                    
* cons
    * 複雑なアプリが対象でないと、overkill 感が強い
    * おそらく切り方が非常に難しい
    
## Reference
* [project-layout](https://github.com/golang-standards/project-layout)
    * github でサンプル
* [how-do-you-structure-your-apps](https://github.com/katzien/talks/blob/master/how-do-you-structure-your-apps/gopherconuk-2018-08-03/slides.pdf)
    * [youtube](https://www.youtube.com/watch?v=oL6JBUk6tj0)
    * [sample](https://github.com/katzien/go-structure-examples)
* [Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
    * Root package is for domain types
    * Group subpackages by dependency
    * Use a shared mock subpackage
    * Main package ties together dependencies
* [Style guideline for Go packages](https://rakyll.org/style-packages/)
    * Organize by (functional)responsibility
* [CODE LIKE THE GO TEAM](https://talks.bjk.fyi/gcru18-best.html#/)
    * naming conventions と package organization
* [Modern Go Application example](https://github.com/sagikazarmark/modern-go-application)
    * config/logging/error/metrics(Prometheus and Jaeger)/health check/graceful restart/daemon/messaging/db
* [How Uber "Go"es](https://speakerdeck.com/lelenanam/how-uber-go-es)
