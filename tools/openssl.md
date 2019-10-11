---
title: openssl
---

## 使い方
``` bash
# 標準コマンド一覧
$ openssl list-standard-commands
# メッセージダイジェストコマンド一覧
$ openssl list-message-digest-commands
# 暗号スイートのコマンド一覧
$ openssl list-cipher-commands

# TLS でサーバに接続
$ openssl s_client -connect www.google.com:443

# self-signed certificate 作成
# -nodes: don't protect private key with a passphrase
# -subj: 対話形式でなく、コマンドラインで指定できる
# /CN=<FQDN>
$ openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```