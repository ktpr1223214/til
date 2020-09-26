---
title: envsubst
---

## 使い方
設定ファイルに環境変数を差し込むことができる。

* 適当な yml ファイルを準備
``` yml 
docker:
    - image: $IMAGE
```

``` bash
# 基本
$ IMAGE=alpine:latest envsubst < sample.yml 

# docker build 時に使う場合は build-time variables 指定
$ export HTTP_PROXY=http://10.20.30.2:1234
$ docker build --build-arg HTTP_PROXY .
```

## Reference
* [envsubst(1) ](https://man7.org/linux/man-pages/man1/envsubst.1.html)
* [Environment substitution with Docker](https://www.robustperception.io/environment-substitution-with-docker)
