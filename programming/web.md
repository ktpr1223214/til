---
title: Web
---

## SQL
* [kamipo TRADITIONALでは防げないINSERT IGNOREという名の化け物Add Star](https://songmu.jp/riji/entry/2015-07-20-insert-ignore.html)
    * レコードがなかったらインサートの話
        * INSERT ON DUPLICATE KEY UPDATE を使う
        * SELECT して無かったら INSERT して、Duplicate エラーをトラップしたら一旦トランザクションを抜けてからまたトランザクションを張って SELECT しなおす
        * DELETE してから INSERT する
        * 素直に例外をあげる

## Session
* サーバ側でもログイン状態を持つという観点もある
    * サーバ側で消せば、クッキーがあろうとログインはできない

* cookie がセットされていない場合
    * http.StatusUnauthorized

## Cookie

## Endpoint
* [Monitoring the health of your application - The upgraded "/ping" route ](https://www.sohamkamani.com/blog/architecture/2018-09-06-application-health-monitoring/)
