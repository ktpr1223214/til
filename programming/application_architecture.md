---
title: Application Architecture
---

## Application Architecture


## 注意
### 依存関係
"依存しない"というのは、パッケージ同士の依存関係が存在しないということだけではなく、そのレイヤが知るべき概念のみレイヤ内で扱うという観点も存在します。

そのため、本来は users テーブルと user_status テーブルに アクセスするのではなく、users テーブルに対する DAO と user_status テーブルに対する DAO を2つ実装し、
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
* [webフロントエンドからwebAPIを呼び出すのをwrapするアレの名前](https://nekogata.hatenablog.com/entry/2019/06/30/211904)
* [Service LocatorとDependency InjectionパターンとDI Container](https://www.nuits.jp/entry/servicelocator-vs-dependencyinjection)
* [UseCaseの再利用性](https://yoskhdia.hatenablog.com/entry/2016/10/18/152624)

### architecture
* [Write code that is easy to delete, not easy to extend.](https://programmingisterrible.com/post/139222674273/write-code-that-is-easy-to-delete-not-easy-to)
* [Modern Software Over-Engineering Mistakes](https://medium.com/@rdsubhas/10-modern-software-engineering-mistakes-bc67fbef4fc8)
* [Ready for changes with Hexagonal Architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749)
  * Netflix

### API design
* [API 設計ガイド](https://cloud.google.com/apis/design/)
  * ネットワーク API の一般的な設計ガイド

### Test
* [Writing a Unit Test](https://developers.mattermost.com/contribute/server/rest-api/#writing-a-unit-test)
    * REST API のテストの観点が書いてあって参考になる

### DDD
* [ボトムアップドメイン駆動設計](https://nrslib.com/bottomup-ddd/)
