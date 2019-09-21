---
title: MySQL
---

## MySQL

## 基本
* system_time_zone
    * システムタイムゾーン
    * ホストマシーンのタイムゾーンで、起動時に自動で設定されその後は変わらない
        * 起動時に明示的に設定したい場合は、環境変数 TZ を mysqld の起動前に設定しておく
* time_zone
    * サーバーの現在のタイムゾーン
    * SYSTEM が初期値で、これはサーバーのタイムゾーンがシステムタイムゾーンと同じということを表す
    * TIMESTAMP カラムには影響(なので、タイムゾーンを意識するならこっち)
        * TIMESTAMP カラムの値は、ストレージでは現在のタイムゾーンから UTC に、読み出しでは UTC から現在のタイムゾーンに変換
        * INSERT: time_zone 指定タイムゾーンの時刻と解釈し、それを UTC に変換して保存
        * SELECT: UTC で保存された時刻を、time_zone 指定タイムゾーンでの表記に変換して表示
    * DATE、TIME、DATETIME カラムの値には影響しない
        * INSERT: そのまま
        * SELECT: そのまま

* TIMESTAMP の format は
    * '1970-01-01 00:00:01' UTC to '2038-01-19 03:14:07' UTC

* timestamp と datetime の違い
    * 

## コマンド
``` bash
# 接続
# ex. mysql -h 127.0.0.1 -P 3306 -u root sample_db -p -e 'select * from sample_table'
$ mysql -h <host> -P <port> -u <user> <db> -p -e 'select * from sample_table'

# 文字コード関連確認
$ show variables like '%char%';

# timezone 関連確認
$ select @@global.system_time_zone, @@global.time_zone, @@session.time_zone;
```

## conf
* my.cnf
``` 
[mysqld]
character-set-server=utf8mb4
collation-server=utf8mb4_general_ci
default-storage-engine=InnoDB
explicit_defaults_for_timestamp=1
default-time-zone=Asia/Tokyo

[client]
loose-default-character-set=utf8mb4
```

* character-set-server    
    * [MySQLの文字コード事情 2017版](https://www.slideshare.net/tmtm/mysql-2017)
* explicit_defaults_for_timestamp
    * timestamp のデフォルト値設定を許可しない

## クエリサンプル
``` bash
# Primary Key の付け替えを oneliner で
$ alter table <table_name> drop primary key, add primary key(<column>); 

```

## docker
* [mysql](https://hub.docker.com/_/mysql/)

## RDS
``` bash
$ sudo yum install mysql
$ mysql -h <endpoint> -P 3306 -u <user_name> -p <db_name>
```

