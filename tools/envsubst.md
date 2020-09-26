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
$ IMAGE=alpine:latest envsubst < sample.yml 
```

## Reference
* [envsubst(1) ](https://man7.org/linux/man-pages/man1/envsubst.1.html)
