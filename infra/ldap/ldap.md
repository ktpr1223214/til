---
title: LDAP
---

## LDAP
* バインド
  * LDAP サービスへログインすること
  * バインドで、LDAP の検索などが出来るように
* バインド DN
  * LDAP サービスへログインするときのユーザ
* ベースDN
  * LDAP サービスへログインした後、どの OU 配下の情報を扱うか

## tools
* ldapsearch
  * ```which ldapsearch```
* ldapadd
  * ```which ldapadd```
``` bash
$ ldapadd -D cn=admin,dc=sample,dc=com -W -f ./sample_user.ldif

# -h: LDAP サーバの IP や FQDN を指定 -x: SASL ではなく簡易認証を使うオプション
# -D: バインド DN を指定 -W: -D で指定したバインド DN に対するバインドパスワードの入力を求める
# -b: ベース DN を指定する
$ ldapsearch -h localhost -x -D "cn=admin,dc=sample,dc=com" -W -b "ou=eigyo,dc=sample,dc=com" cn=test\*
$ ldapsearch -h localhost -x -D "cn=admin,dc=sample,dc=com" -W -b "cn=admin,dc=sample,dc=com" cn=admin\*
```
* Apache Directory Studio
  * [Downloads](https://directory.apache.org/studio/downloads.html)