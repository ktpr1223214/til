---
title: Go
---

## Go

## Goroutine
### Scheduler
* [The Go scheduler](https://morsmachine.dk/go-scheduler)
  * 全体感わかりやすい
* OS が既にスケジューラーを持つのに、なぜ userspace スケジューラーを Go は持つのか
  * OS thread の overhead(実質 process なわけで)
  * OS は、Go のモデル(メモリとかの話？)を考慮したスケジューリングが出来ない
    * ex. GC の時
#### [Thread のモデルは大まかに3種](https://en.wikipedia.org/wiki/Thread_(computing)#N:1_(user-level_threading))
  * N:1
    * An N:1 model implies that all application-level threads map to one kernel-level scheduled entity;the kernel has no knowledge of the application threads
      * pros: context switch が高速(あくまでアプリ内でのユーザスレッドの切替としての context switch なわけで、OS を介さないので）
      * cons: multi-core を有効利用できず、ハードレベルでの高速化は利用できない
        * ex. あるスレッドが、I/O request を処理している時、プロセス自体が block されてしまう(kernel じゃないんだしそりゃそう)
  * 1:1
    * Threads created by the user in a 1:1 correspondence with schedulable entities in the kernel
      * pros: 全 core を有効活用できる
      * cons: context switch が、OS(ex. Linux Kernel)を介すので、遅い
  * M:N

* [Go: Goroutine, OS Thread and CPU Management](https://medium.com/a-journey-with-go/go-goroutine-os-thread-and-cpu-management-2f5a5eaf518a)
    * 基本的な話がまとまっている

#### The main concepts
* G
  * goroutine
  * G は stack や instruction pointer などを持つ。また、ブロックされる可能性があるチャネルなど、スケジューリングに重要な情報も持つ
    * 2KB stack
  * P の local runqueue or global runqueue に配置
* M
  * worker thread, or machine.
* P
  * processor, a resource that is required to execute Go code.
    M must have an associated P to execute Go code, however it can be
    blocked or in a syscall w/o(without) an associated P.
  * context for scheduling. You can look at it as a localized version of the scheduler which runs Go code on a single thread.
    * N:1 から M:N scheduler になるのに重要
  * local runqueue を持つ

M が OS Thread で、P が論理 CPU
Goroutine を実行するには、M は P を持つ必要がある
P の数は GOMAXPROCS で設定され、GOMAXPROCS 数がある時点で実行される Go code の最大数（でいいよね？）
一方で、M 自体はそれよりも多く生成することができる

global runqueue も存在
* The global runqueue is a runqueue that contexts pull from when they run out of their local runqueue.
* Contexts also periodically check the global runqueue for goroutines.

Goroutines are added to the end of a runqueue whenever a goroutine executes a go statement. Once a context has run a goroutine until a scheduling point, it pops a goroutine off its runqueue, sets stack and instruction pointer and begins running the goroutine.

* なぜ context が必要なのか
  * thread が実行中にブロックされた場合、context を別の thread に渡して実行を続けることができる
  * thread に runqueue が直接ひっつく形だと、そういうことは出来ない

## Stack
* [Contiguous stacks in Go](https://agis.io/post/contiguous-stacks-golang/)
  * わかりやすい説明
* [Design doc](https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub)

## Reference
### 公式など
* [go](https://github.com/golang/go/wiki)
    * 公式
* [Go Proverbs](https://go-proverbs.github.io/)
* [Russ cox](https://research.swtch.com/)
* [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)

### package
* [Go Sub-repository Packages](https://godoc.org/-/subrepo)
    * These packages are part of the Go Project but outside the main Go tree. They are developed under looser compatibility requirements than the Go core.
    * ココらへんの機能は変に自作せずに使うと良さげ

### Go 開発全般
* [Practical Go: Real world advice for writing maintainable Go programs](https://dave.cheney.net/practical-go/presentations/qcon-china.html)
    * Dave Cheney
* [go-for-industrial-programming](https://peter.bourgon.org/go-for-industrial-programming/)
* [how-to-ship-production-grade-go](https://www.oreilly.com/ideas/how-to-ship-production-grade-go)
    * Wrap errors
    * Report panics
    * Use structured logs
    * Ship application metrics
    * Write more tests than you think you should
* [idiomatic-go](https://about.sourcegraph.com/go/idiomatic-go/)
    * べからず集として良い
* [gotraining](https://github.com/ardanlabs/gotraining)
    * 色々なリンク類。Go + alpha でかなり網羅的
* [Good Code vs Bad Code in Golang](https://medium.com/@teivah/good-code-vs-bad-code-in-golang-84cb3c5da49d)

### 歴史や背景・哲学
* [Is Go An Object Oriented Language?](https://spf13.com/post/is-go-object-oriented/)

### Go 文法・内部詳細
* [go-internals](https://github.com/teh-cmc/go-internals)
* [The Go Memory Model](https://golang.org/ref/mem)
* [The Go Object Lifecycle](https://middlemost.com/object-lifecycle/)
* [Go Interfaces](https://www.airs.com/blog/archives/277)
* [Understanding Nil](https://speakerdeck.com/campoy/understanding-nil)
* [Go Walkthrough](https://medium.com/go-walkthrough)
    * io など package について解説
* [Go memory ballast: How I learnt to stop worrying and love the heap](https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap-26c2462549a2/)
    * GC の話なども書いてある

## Podcast
* [gotime](https://changelog.com/gotime)
