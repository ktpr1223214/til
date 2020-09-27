# graceful-shutdown
``` bash
# time.Sleep(5 * time.Second) の場合
$ curl -v http://localhost:8888/health
# を実行してから、server 側 ctrl+C でこの curl を返すまでは待つ

# time.Sleep(20 * time.Second) の場合
$ curl -v http://localhost:8888/health
# を実行してから、server 側 ctrl+C でこの curl はレスポンスが返らない
# エラー: could not gracefully shutdown the server: context deadline exceeded も出力される
```

## Reference
* [gorilla/mux](https://github.com/gorilla/mux#graceful-shutdown)
* [GoでGraceful Shutdown](https://christina04.hatenablog.com/entry/go-graceful-shutdown)
