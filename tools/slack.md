---
title: Slack
---

## Slack

## App
* [このページ](https://api.slack.com/apps)から、新規作成
* 権限は細かく制御できるので、必要なものをつける
  * Permissions -> Scopes から
* bot が作りたいなら、
  * Features の項目から Bot user を追加
  * Interactive Components を追加。ここで Request URL が必要
    * https://api.slack.com/tutorials/tunneling-with-ngrok
* とりあえずの設定が済んだら、Install App することで、token などが発行される

## tips
* channel id の取得
  * URL copy したときに、https://~/messages/<channel id> という構成で得られる
* private channel への bot からの post
  * invite しておかないと、channel_not_found でエラーになるはず
