---
title: Delve
---

## Delve
デバッガー

## How to
``` bash
# install(GO111MODULE=on が有効なときは、一時的に off にして global にインストール)
$ GO111MODULE=off go get -u github.com/go-delve/delve/cmd/dlv
# or

# help 
$ dlv help

# run
$ dlv debug ~/main.go

# set breakpoint
$ b main.main

# skip to breakpoint
$ c

# print variable
$ p <variable_name>

# next
$ n 

# stepin
$ s

# stepout
$ stepout

# exit
$ q
```