---
title: io
---

## io
標準パッケージ
* byte stream を扱うために役立つ
    
## Reader
``` go
type Reader interface {
        Read(p []byte) (n int, err error)
}
```
* buffer p を渡す構造なので、再利用が可能