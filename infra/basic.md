---
title: Infra basic
---

## 基礎知識
* cpu
    * utilization
        * 使っている or not の 0/1 なので、その一定時間の平均を utilization としている

### コマンド
* 

* ps
    
``` 
# -a: 端末操作のプロセスも表示
# -x: 現在実行しているプロセスを表示
# -u: CPU やメモリの使用率表示
# -f: プロセスを階層で表示
$ ps -auxf

```

* netstat
    * 
``` bash
$ 
```

* dstat
    * https://qiita.com/harukasan/items/b18e484662943d834901

* https://qiita.com/kunihirotanaka/items/70d43d48757aea79de2d#%E3%83%90%E3%83%83%E3%83%95%E3%82%A1%E3%82%AD%E3%83%A3%E3%83%83%E3%82%B7%E3%83%A5%E3%81%AE%E6%8C%99%E5%8B%95%E3%82%92%E7%A2%BA%E8%AA%8D%E3%81%99%E3%82%8B

```
# -t: 時間つける
# -c(m/d/n): それぞれ、CPU 使用率・メモリ使用量・Disk IO・ネットワーク IO
# -s: スワップの used/free
# usr: ユーザプロセス使用率 sys: システムプロセス使用率 idl: 空き率 wai: プロセスの待ち状態率 hiq: ハードウェアの割り込み率 siq: ソフトウェアの割り込み率
# used: メモリ使用サイズ buff: バッファキャッシュサイズ(バッファキャッシュ) cach: キャッシュサイズ(ページキャッシュ) free: 未使用メモリサイズ
# バッファキャッシュはブロックデバイス(ex. ハードディスク)を直接アクセスするときに使用されるキャッシュ
# ページキャッシュはファイルシステムに対するキャッシュであり、ファイル単位でアクセスするときに使用されるキャッシュ
# used: 使用サイズ free: 空きサイズ
# 
# read: ディスク読み込みのバイト数 writ: ディスク書き込みのバイト数
# recv: 受信量 send: 送信料
$ dstat -tcmsdn

# --tcp: enable tcp stats (listen, established, syn, time_wait, close)
# 数の単位はソケット数かと
# listen: コネクション待受数 established: コネクションが開かれ、データ転送が行える状態の数
# syn: TCP の 応答確認とかを行っている数
# time_wait: コネクション終了要求応答確認をリモートホストが確実に受取るのに必要な時間が経過するまで待機している数
# close: コネクションが存在せず、待受でもない状態の数
$ dstat --tcp
```

### ssh
* authorized_keys
    * サーバ側で、接続を許可する公開鍵を登録するファイル

### FAQ
* sudo を付けたときに sudo: unable to resolve host が表示された場合
``` bash
# /etc/hosts を確認
$ cat /etc/hosts
# /etc/hosts に host 名を追加
$ sudo sh -c 'echo 127.0.1.1 $(hostname) >> /etc/hosts'
```