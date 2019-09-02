---
title: Bytes
---

## Bytes

## bytes/string 
* []byte: Byte slices 
    * mutable
    * resizable
    * contiguous list of bytes
    
``` go
buf := []byte{1,2,3,4}

// mutable
buf[3] = 5

// resizable 
buf = buf[:2] 

// it’s contiguous so each byte exists one after another in memory
```

* strings:
    * immutable
    * fixed-size
    * contiguous list of bytes    
* パフォーマンス観点からは、string は常に生成が必要となるので、GC に負荷がかかる
* 一方で、開発者には理解しやすい

``` go
func NewReader(b []byte) *Reader
func NewReader(s string) *Reader
``` 
* メモリ上の byte slice or string を wrap した io.Reader を返す
* 更に、io.ReaderAt, io.WriterTo, io.ByteReader, io.ByteScanner, io.RuneReader, io.RuneScanner, io.Seeker も満たしている 
      