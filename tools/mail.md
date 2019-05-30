---
title: Mail
---

## Mail

``` bash
# docker で postfix を実行
$ docker run -d -p 25:25 --name <container_name> -e maildomain=<domain> -e smtp_user=<username>:<password> catatnight/postfix

# Plain Auth 用の
# ex. printf "%s\0%s\0%s" user@example.jp user@example.jp hoge | openssl base64 -e | tr -d '\n'; echo
$ printf "%s\0%s\0%s" <username> <username> <password> | openssl base64 -e | tr -d '\n'; echo

# telnet で実行
$ EHLO localhost

# 認証
$ AUTH PLAIN <上で生成した認証文字列>

# 送信元設定
$ MAIL FROM:gufa@hoge.com
# 送信先設定
$ RCPT TO: ~@~.com

# DATA を書いてその後本文入力(終わりは、. を入力)
$ DATA
# 終了
$ QUIT

# 入る
$ docker exec -it <container_name> bash
$ cat /var/log/mail.log
```