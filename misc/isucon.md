---
title: ISUCON
---

## ISUCON

### tools
* netdata
* alp
* pt-query-digest
* mysql workbench

### 初期セットアップ関連
* コードを git 管理
    * app のコード
    * /etc 配下の設定ファイルも別に git で
* DB データのバックアップ

* ベンチマーク実行
    * ベンチマーク実行前の件数は取得し忘れたが、実行前後で見比べてレコード数に変化があるテーブルと変化の無いテーブルを確認する。変化の無いテーブルはキャッシュするなどの対象になるが今回はキャッシュはしなかった。

* sql は機械的に書き出すも良い
    
* kataribe で endpoint ごとのプロファイリング準備
* netdata によるOSリソースの取得と可視化
* pt-querydigest による mariadb の SQL のプロファイリング準備
* deploy と計測準備用の script を配置等々

### あり得る手段
* OS の設定チューニング
    tcpのlife cycleを短くする
    エフェメラルポートの利用範囲を広げる
    利用できるファイルディスクリプタ数を上げる

### 終わりにやること
* 各種ログの出力を OFF に
* 使ってない MW は全部止める
* GoのTemplateエンジンをやめてコードに埋め込む 

## 事前準備が必要そうなこと
* [ ] 各種再起動の雛形作成
* [ ] デプロイスクリプト
* [ ] /initialize が叩かれたときに実行する、ログなどの退避用スクリプト
* [ ] ssh とかデプロイ周りの権限確認
* [ ] 

## Monitoring

### netdata
``` bash
# install
$ bash <(curl -Ss https://my-netdata.io/kickstart.sh)
$ sudo systemctl start netdata

# docker
$ docker run -d --name=netdata -p 19999:19999 -v /proc:/host/proc:ro -v /sys:/host/sys:ro -v /var/run/docker.sock:/var/run/docker.sock:ro --cap-add SYS_PTRACE --security-opt apparmor=unconfined netdata/netdata
```


## Mysql
* https://dev.mysql.com/downloads/workbench/ も便利
    * ER 図の dump

``` bash
# 
$ use <database>;
$ show tables;
# 全テーブルの情報
$ select table_name, engine, table_rows, avg_row_length, floor((data_length+index_length)/1024/1024) as allMB, floor((data_length)/1024/1024) as dMB, floor((index_length)/1024/1024) as iMB from information_schema.tables where table_schema=database() order by (data_length+index_length) desc;

# あるテーブルの構造
$ describe <table>;
``` 

### pt-query-digest
```
# install(macos)
$ brew install percona-toolkit
 
# https://www.percona.com/downloads/percona-toolkit/LATEST/

# use
$ pt-query-digest /var/log/mysql/slow.log
```

* slow-query log の出力
    * long_query_time = 0で全てのクエリを出力
``` 
[mysqld]
slow_query_log
slow_query_log_file = /var/log/slow.log
long_query_time = 0
```

## alp
``` bash
# install
$ wget https://github.com/tkuchiki/alp/releases/download/v0.4.0/alp_darwin_amd64.zip
$ unzip alp_darwin_amd64.zip
$ mv alp /usr/local/bin/

# use
$ alp -f access.log
```

## chrome
* キャッシュの消去とハード再読み込み

"Memory Cache" stores and loads resources to and from Memory (RAM). So this is much faster but it is non-persistent. Content is available until you close the Browser.

"Disk Cache" is persistent. Cached resources are stored and loaded to and from disk.

## cache
https://developers.cyberagent.co.jp/blog/archives/5975/


## Examples
###
```  
echo: http: Accept error: accept tcp [::]:3000: accept: too many open files; retrying in 5ms
echo: http: Accept error: accept tcp [::]:3000: accept: too many open files; retrying in 10ms
{"time":"2019-07-15T01:54:24.918781+09:00","level":"-","prefix":"-","file":"login.go","line":"80","message":"dial tcp :6379: socket: too many open files"}
```

## Reference
### 全体
* [](https://blog.yuuk.io/entry/web-operations-isucon)
* [自分のチームのISUCONでの戦い方](https://medium.com/@catatsuy/%E8%87%AA%E5%88%86%E3%81%AE%E3%83%81%E3%83%BC%E3%83%A0%E3%81%AEisucon%E3%81%A7%E3%81%AE%E6%88%A6%E3%81%84%E6%96%B9-c8fe121316aa)
    * [スクショ・各パスの役割・どういうアプリケーションか・キーになる関数 catatsuy/isucon7-qualifier](https://github.com/catatsuy/isucon7-qualifier/issues/1)
    * [こっちのが updated?](https://gist.github.com/catatsuy/e627aaf118fbe001f2e7c665fda48146)
* https://github.com/wantedly/shisucon2019-teppei-haruki/wiki#%E6%99%82%E9%96%93%E9%85%8D%E5%88%86
    
### db
* [](http://dsas.blog.klab.org/archives/2018-02/configure-sql-db.html)