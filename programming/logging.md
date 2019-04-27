---
title: Logging
---

## Logging

### Structure Logging
* pros
    * 機械的な処理が容易
    * 形式を明示的に指定可能
* cons
    * シリアライズにコスト

* 形式の明示的な指定の例
``` go
// userId, requestId, endpoint をこの logger を使うすべての statement に付与 
log := logger.With(zap.String("userId", "fuga"), zap.String("requestId", "piyo"))
log.Info("hello world")
```

* ログ量が少ない・人しか読まない・grep/tail で出来るより複雑なことはしない、であれば別に非構造化ログでも良いがそうでないなら構造化ログ推奨
    * 入門監視

## Reference
* [Logging best practices to get the most out of application level logging – Slides](https://geshan.com.np/blog/2019/03/follow-these-logging-best-practices-to-get-the-most-out-of-application-level-logging-slides/)
    * [slide](https://speakerdeck.com/geshan/logging-best-practices)
