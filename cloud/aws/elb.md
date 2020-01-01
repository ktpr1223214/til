---
title: ELB
---

## ELB
作成すると、DNS 名が振られるので例えば以下で ip を得ることが出来る
* 勿論 SG で接続を許可している場合
``` 
# ip を返す
$ dig <DNS name>
```

## CNAME
* 上述の通り、DNS 名が振られて A レコードが存在している
* なので、更に CNAME で別名を振ることも可能
    * 例えば ECS の host based routing などもここから出来る

## NLB 
* NLB と Proxy Protocol
  * [NLB (Network Load Balancer) が Proxy Protocol に対応しました](https://dev.classmethod.jp/cloud/aws/nlb-meets-proxy-protocol-v2/)
