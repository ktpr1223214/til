---
title: Curl
---

## Curl

## 使い方
``` bash
# 詳細(-v)
$ curl -v ~

# Get(-G/--get) --data-urlencode は -G と一緒の場合、URL の末尾にクエリーを付与
$ curl -G --data-urlencode "search word" http://localhost:8888

# ヘッダー行
$ curl -H "Fuga: Hoge" http://localhost:8888

# リクエストを送る(--request==<method> or -X <method>)
$ curl -X POST http://localhost:8888

# ボディの送付(-d)
# デフォルトでの、Content-Type は application/x-www-form-urlencoded
$ curl -d "{\"hello\": \"world\"}" -H "Content-Type: application/json" http://localhost:8888

# -d で複数送信 title=The Art of Community&author=Jono Bacon というボディとなる(各項目が & でつながれる)
$ curl -d title="The Art of Community" -d author="Jono Bacon" http://localhost:8888

# URL エンコードで送付 title=The%20Art%20of%20Community&author=Jono%20Bacon というボディになる
$ curl --data-urlencode title="The Art of Community" --data-urlencode author="Jono Bacon" http://localhost:8888

# Accept-Encoding: deflate, gzip ヘッダーを送付(--compressed)
$ curl --compressed http://localhost:8888

# -c/--cookie-jar で指定したファイルに受信したクッキーを保存 -b/--cookie で指定したファイルから読み込んでクッキー送信
# -b は "name=value" のように、個別の値を指定することも可能
# ブラウザのように送受信を行いたい場合には、2つとも指定
$ curl -c cookie.txt -b cookie.txt http://localhost:8888

# -u/--user ユーザ名とパスワードを送信
$ curl -u user:pass http://localhost:8888

# プロキシ -x/--proxy プロキシ認証 -U/--proxy-user

# レスポンス300台かつレスポンスヘッダーに Location ヘッダーがあった場合、そのヘッダーで指定された URL に再度リクエスト(-L)
$ curl -L http://localhost:8888

# 複数のリクエストを並べることで、Keep-Alive 利用
$ curl -v http://localhost:8888 http://localhost:8888

# TLS 系
# TLS で接続(-1/--tlsv1)
$ curl -1 http://localhost:8888

# chunk
# -T でファイル転送とともに、chunked
$ curl -T README.md -H "Transfer-Encoding: chunked" http://localhost:8888

# ファイル保存(-O/--remote-name)
$ curl -O http://example.com/download/sample.pdf
```