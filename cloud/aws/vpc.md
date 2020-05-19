---
title: VPC
---

## Route table
* VPC もルートテーブルを持つ
  * default VPC の場合、インターネット接続を持つメインルートテーブルが紐づいている
    * デフォルトでは、デフォルト以外の VPC を作成すると、メインルートテーブルにはローカルルートのみ
  * subnet は特定のルートテーブルに明示的に関連付けることができるが、それ以外の場合、サブネットはメインルートテーブルに暗黙的に関連付けられる

## VPC peering connection
* VPC peering connection 設定後、互いの route table が、peering connection を向くようにする必要がある
  * でないと、双方向での通信は疎通できない
