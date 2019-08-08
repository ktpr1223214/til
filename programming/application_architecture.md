---
title: Application Architecture
---

## Application Architecture


## 注意
### 依存関係
"依存しない"というのは、パッケージ同士の依存関係が存在しないということだけではなく、そのレイヤが知るべき概念のみレイヤ内で扱うという観点も存在します。
そのため、本来は users テーブルと user_status テーブルに アクセスするのではなく、
users テーブルに対する DAO と user_status テーブルに対する DAO を2つ実装し、
モデルレイヤにてそれらを利用して domain.User を作成するべきでしょう。 

### トランザクションをどこで扱うべきか

## Reference
* [Thoughts on Code Organization](https://medium.com/@egonelbre/thoughts-on-code-organization-c668e7cc4b96)
    * When the code-base is less than 30,000 lines of code, then code organization isn’t a big hindrance in getting work done.
        * This number also suggests that you cannot draw any significant conclusions from code-bases that have less than 30K LOC (lines of code).
    * Just organizing things doesn't necessarily create value
* [Benefits of dependencies in software projects as a function of effort](https://eli.thegreenplace.net/2017/benefits-of-dependencies-in-software-projects-as-a-function-of-effort/)    
* [ドメイン駆動設計解説シリーズ](https://little-hands.hatenablog.com/entry/top)
* [レイヤードアーキテクチャを振り返る](https://buildersbox.corp-sansan.com/entry/2019/04/21/000000_1)