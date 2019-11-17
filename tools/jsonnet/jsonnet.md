---
title: Jsonnet
---

## setup
* install jsonnet etc
``` bash
$ brew install jsonnet jq
$ go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
$ jb init
```

* setup
``` bash
# example
$ jsonnet -J . main.jsonnet
```
