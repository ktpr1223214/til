---
title: CloudMapper
---

## how to
* Docker で動かす手順
* 事前: 
  * aws-vault を入れて、設定もしておく
  * config.json を config.json.demo を参考に設定

``` bash
$ docker build -t cloudmapper .
$ aws-vault exec <profile> --server --
$ docker run -p 8000:8000 -it cloudmapper /bin/bash

# コンテナで
$ pipenv shell

# network
$ python cloudmapper.py prepare --account <account_name>
$ python cloudmapper.py webserver --public

# 
$ python cloudmapper.py collect --account <account_name>
# 
$ python cloudmapper.py report --account <account_name>
```
