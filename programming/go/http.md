---
title: HTTP
---

## HTTP
* go の主に標準ライブラリ中心に

### Handler
* Handler
    * ServeHTTP(http.ResponseWriter, *http.Request) というメソッドを持つインターフェース

### Request
* 構成要素
    * URL
    * Header
    * Body
    * Form/PostForm/MultipartForm

* URL の形式
    * scheme://[userinfo@]host/path[?query][#fragment]
        * fragment はブラウザからのリクエストでは取得できない
            * サーバに送信される前にブラウザに除去されるので
        * HTTP クライアントライブラリや、クライアントフレームワークからリクエストを受け取る場合に有効

* Body
    * Request/Response いずれの場合もボディは Body フィールド
    * io.ReadCloser インターフェース

* フォームフィールド群
    * Form
        * 呼ぶ必要のあるメソッド: ParseForm
        * キーバリューペアの取得元: URL/フォーム
        * コンテントタイプ: URL エンコード
        * URL/フォーム両方からキーバリューペアを取り出せる。URL とフォームに同一キーが存在する場合、フォームを優先して先に配置
    * PostForm
        * 呼ぶ必要のあるメソッド: Form(? ParseForm ではなくて？)
        * キーバリューペアの取得元: フォーム
        * コンテントタイプ: URL エンコード
        * URL は無視でフォームのキーバリューペアのみ
        * application/x-www-form-urlencoded のみしかサポートされていないため、マルチパートのキーバリューペアを取得するには、MultipartForm フィールドが必要
    * MultipartForm:
        * 呼ぶ必要のあるメソッド: ParseMultipartForm
        * キーバリューペアの取得元: フォーム
        * コンテントタイプ: マルチパート
    * FormValue
        * 呼ぶ必要のあるメソッド: なし
        * キーバリューペアの取得元: URL/フォーム
        * コンテントタイプ: URL エンコード
    * PostFormValue
        * 呼ぶ必要のあるメソッド: なし
        * キーバリューペアの取得元: フォーム
        * コンテントタイプ: URL エンコード

### ResponseWriter
* Write
    * バイト配列を受け取り、HTTP レスポンスのボディに書き込み
    * 呼び出しまでにヘッダでコンテンツタイプが設定されていない場合は、データの先頭512バイトでコンテンツタイプ判定

* WriteHeader
    * HTTP レスポンスのステータスコードを引数に、返すステータスコードの書き込み
    * このメソッド呼び出し以降は、ヘッダに書き込むことはできなくなる
    * このメソッドを呼び出さない場合の Write のデフォルトは200

* Header
    * 変更可能なヘッダのマップを返す
    * 変更されたヘッダはクライアントに送信される HTTP レスポンスに入る

### Cookie
* 有効期限の指定
    * Expires or MaxAge フィールド
        * Expires: いつそのクッキーが期限切れになるか
        * MaxAge: ブラウザ内でそのクッキーが生成されてからどれだけの期間（秒）有効か
    * Expires は HTTP1.1 で MaxAge 優先で非推奨となったがほぼ全てのブラウザで対応
    * MaxAge は、IE6・7・8で未対応
    * なので、現実的には Expires のみ or 両方使う
* Expires フィールドが設定されていない場合、そのクッキーはセッションクッキー(ブラウザ閉じられると削除)で、それ以外は永続性クッキー(期限切れ or 除去されない限り維持)
* クッキーの値は URL エンコードする必要がある
    * go では ```base64.URLEncoding.EncodeToString(~)```

### その他
* [Why Handler registered by http.HandleFunc is called twice?](https://groups.google.com/forum/#!topic/golang-nuts/1sgaQGpIILM)

## Reference
* [So you want to expose Go on the Internet](https://blog.cloudflare.com/exposing-go-on-the-internet/)
  * Applying timeouts is a matter of resource control. Even if goroutines are cheap, file descriptors are always limited
  * A zero/default http.Server, like the one used by the package-level helpers http.ListenAndServe and http.ListenAndServeTLS, comes with no timeouts -> 避ける
  * IdleTimeout: HTTP の Keep-Alive の idle time 設定
  * TCP Keep-Alives
    * こっちは TCP layer（かつ OS レベル）
    * https://github.com/golang/go/blob/38543c2813a1075e09693894625421309d8ef333/src/net/tcpsockopt_unix.go#L15
      * 設定はこんな感じ
* [The complete guide to Go net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
  * server/client 側両方について
  * SetDeadline: TCP/UDP のレイヤ（transport layer）のレベルの機能（net.Conn）
