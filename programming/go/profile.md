---
title: Profile
---

## Profile

### pprof
* CLI といった短時間実行のプログラムと、Web アプリのように常駐プログラムでプロファイル方法が少し変わる
    * 前者には、[pkg/profile](https://github.com/pkg/profile) を使うと吉
    * 

* 作業手順は、
    1. コードにプロファイラを埋め込む
    2. コードを実行してプロファイル実行
    3. プロファイル結果の確認

go tool pprof -seconds 120

``` bash
# web gui 使った場合もプロファイル結果ファイルは保存されているはずなので、それを指定してインタラクティブ実行可能
# ex. 
$ go tool pprof /Users/um003404/pprof/pprof.samples.cpu.003.pb.gz

# interactive mode
# 実行時間が多い順(数字指定で上からいくつみるかを指定)
$ top

# -cum: その関数自体だけでなく、その関数から呼びだされた関数の実行時間も合わせた合計時間でソート
$ top -cum
```

## Reference
* [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html)