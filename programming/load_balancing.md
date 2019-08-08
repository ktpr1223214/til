---
title: Load Balancing
---

## Load Balancing
### 概要
* Client と Backends の間でいくつかの重要なタスクを担う
    * Service discovery
        * What backends are available in the system? What are their addresses?
    * Health checking
        * What backends are currently healthy and available to accept requests?
    * Load balancing
        * What algorithm should be used to balance individual requests across the healthy backends?

### Benefits
* Naming abstraction
* Fault torelance
* Cost and performance benefits

### L4 と L7
* LB は大きくどちらかに分類される
    * 注意: あまりこの言葉に引っ張られない方が良いらしい

### Load Balancer の機能
LB に何を用いるかで勿論対応状況は変わるので注意(あくまで一般的な機能リスト)


## Reference
* [Introduction to modern network load balancing and proxying](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236)
    * 網羅的に纏まっていて良い
