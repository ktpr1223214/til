---
title: Go
---

## Go

## Goroutine
### Scheduler
* [HACKING.md](https://github.com/golang/go/blob/master/src/runtime/HACKING.md)
  * 公式 runtime の programming document
* [proc.go#L19](https://github.com/golang/go/blob/20a838ab94178c55bc4dc23ddc332fce8545a493/src/runtime/proc.go#L19)
  * コード中の説明
* [go-runtime-scheduler](https://speakerdeck.com/retervision/go-runtime-scheduler)
  * slide
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
  * [G struct](https://github.com/golang/go/blob/5dd978a283ca445f8b5f255773b3904497365b61/src/runtime/runtime2.go#L332)
* M
  * worker thread, or machine.
  * [M struct](https://github.com/golang/go/blob/5dd978a283ca445f8b5f255773b3904497365b61/src/runtime/runtime2.go#L403)
* P
  * processor, a resource that is required to execute Go code.
    M must have an associated P to execute Go code, however it can be
    blocked or in a syscall w/o(without) an associated P.
  * context for scheduling. You can look at it as a localized version of the scheduler which runs Go code on a single thread.
    * N:1 から M:N scheduler になるのに重要
  * resources required to execute user Go code, such as scheduler and memory allocator state
    * A P can be thought of like a CPU in the OS scheduler and the contents of the p type like per-CPU state. This is a good place to put state that needs to be sharded for efficiency, but doesn't need to be per-thread or per-goroutine.
  * local runqueue を持つ
  * [P struct](https://github.com/golang/go/blob/5dd978a283ca445f8b5f255773b3904497365b61/src/runtime/runtime2.go#L470)

M が OS Thread で、P が論理 CPU
Goroutine を実行するには、M は P に assign される必要がある。
P 毎に一度に実行される M の数は1つまで（当然）。
P の数は GOMAXPROCS で設定され、GOMAXPROCS 数がある時点で実行さ
れる Go code の最大数（でいいよね？）
つまり、user level Go code を同時に実行する thread の
最大数が GOMAXPROCS
一方で、M 自体はそれよりも多く生成することができる。
つまり、Go code のための system call 呼び出しで block される
thread の数に制限はなく、GOMAXPROCS の数にカウントしない。

global runqueue も存在
* The global runqueue is a runqueue that contexts pull from when they run out of their local runqueue.
* Contexts also periodically check the global runqueue for goroutines.

Goroutines are added to the end of a runqueue whenever a goroutine executes a go statement. Once a context has run a goroutine until a scheduling point, it pops a goroutine off its runqueue, sets stack and instruction pointer and begins running the goroutine.

* なぜ context が必要なのか
  * thread が実行中にブロックされた場合、context を別の thread に渡して実行を続けることができる
  * thread に runqueue が直接ひっつく形だと、そういうことは出来ない

#### 大雑把な概要
* [GO SCHEDULER: MS, PS & GS](https://povilasv.me/go-scheduler/)
基本的な動きとしては、OS が thread を実行しそこで Go のコードも実行される。
Go 側の挙動としては、Go compiler が Go runtime への呼び出しを
仕込んでおくことにより、scheduler に通知・アクションを実行することができる。

runtime への呼び出し(scheduling decision)を行う可能性のある操作として
* The use of the keyword go
* Garbage collection
* System calls
* Synchronization and Orchestration
  * atomic、mutex、channel 操作など

* [Go's work-stealing scheduler](https://rakyll.org/scheduler/)
  * 絵が分かりやすくて良い

scheduling 処理の1round
```
runtime.schedule() {
    // only 1/61 of the time, check the global runnable queue for a G.
    // if not found, check the local queue.
    // if not found,
    //     try to steal from other Ps.
    //     if not, check the global runnable queue.
    //     if not found, poll network.
}
```

#### spinning thread
* [proc.go#L31](https://github.com/golang/go/blob/go1.14/src/runtime/proc.go#L31)
  * コード中の説明
* [Design doc: Syscalls/M Parking and Unparking](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit)

### netpoller
* [The Go netpoller](http://morsmachine.dk/netpoller)
  * The part that converts asynchronous I/O into blocking I/O is called the netpoller. It sits in its own thread, receiving events from goroutines wishing to do network I/O. The netpoller uses whichever interface the OS provides to do polling of network sockets.
* [Scheduling In Go : Part II - Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)
  * William Kennedy のやつ
  * By using the network poller for networking system calls, the scheduler can prevent Goroutines from blocking the M when those system calls are made. This helps to keep the M available to execute other Goroutines in the P’s LRQ without the need to create new Ms. This helps to reduce scheduling load on the OS.
    * 要は network IO について、M(OS thread) を無駄に増やす必要がなくなるということかと
  * The big win here is that, to execute network system calls, no extra Ms are needed. The network poller has an OS Thread and it is handling an efficient event loop.

## Stack
* [Design doc](https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub)
* [Contiguous stacks in Go](https://agis.io/post/contiguous-stacks-golang/)
  * わかりやすい説明
* [Go: How Does the Goroutine Stack Size Evolve?](https://medium.com/a-journey-with-go/go-how-does-the-goroutine-stack-size-evolve-447fc02085e5)

## Memory


## GC
* [Garbage Collection In Go : Part I - Semantics](https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html)
  * William Kennedy の Series


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
