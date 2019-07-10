# bytes
Go の bytes/string 辺りを調査

## 検証環境
* macOS 10.12.2(Sierra) 
* go version go1.12.5 darwin/amd64

## 検証
``` bash
# escape analysis
$ go build -gcflags '-m' main.go

# benchmark
$ go test -bench . -benchmem
```
