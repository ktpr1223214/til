---
title: HTTP
---

## cache
* [HTTP キャッシュ](https://developer.mozilla.org/ja/docs/Web/HTTP/Caching)


## 注意
* http://localhost:8000/ とかにアクセスすると、以下の事象で2回アクセスが発生するので注意！！！！！
    * https://github.com/golang/go/issues/1298
``` 
forgive my careless, GO works well, it's my fault. I used Chrome Browser to test the
program, when i open the url(http://localhost:12345/), Chrome will send 2 requests, one
is http://localhost:12345, another is http://localhost:12345/favicon.ico, and i read the
document about http package, it's said both the 2 requests will be handled, that's why i
get two outputs every time
```

## Reference
* [httpbin.org](http://httpbin.org/)