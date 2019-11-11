---
title: jq
---

## jq
``` bash
# sample.json を 以下で定義しておく
{
  "key": "value"
}
# この時、以下のコマンドを実行すると、以下の結果が出力される
$ cat sample.json | jq '{dashboard: .} * {dashboard: {sample: true}}'
{
  "dashboard": {
    "key": "value",
    "sample": true
  }
}
# 挙動としては、{dashboard: .} の . に sample.json の内容が入り、更に後ろの JSON も追加される

# この結果からからもわかるかと 
$ cat sample.json | jq '. * {sample: true}'
{
  "key": "value",
  "sample": true
}
```
