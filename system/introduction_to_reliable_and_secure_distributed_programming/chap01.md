---
title: Introduction
---

## Motivation
* 分散アルゴリズムは何かしらプロセス群が協調して動作するアルゴリズムに関する分野
* それらプロセスが単に並行して稼働するのみならず、クラッシュや分断が発生するプロセスも存在することを許容するのがポイント！
    * これを partial failures といい、分散システムの大きな特徴
* partial failures が発生したとしても、残ったプロセスがうまく協調を続けなければならない。そういったロバストなシステムなりアルゴリズムを考えたい
    * で、色々難しいところがあるわけで、例えばクラッシュって簡単に書くけど本当にクラッシュしたかどうかを検知するのも難しい
        * 例えば ping で確認するとした場合、レスポンスが無いのは単にネットワークの問題なのかもしれない、とか
* こういう複雑な現実をそのままで考えるのは難しいので、抽象化したモデルでアルゴリズムを考える・性質を示す、といったことが分散システムの理論でやりたいこと

以下で、分散システムのモデルの表現の構成要素について見ていく。
chapter2.1までの中心的な話題をここに含める。

## 分散システムの抽象化
この本では分散アルゴリズムを asynchronous event-based composition model で表現。
雑にイメージを書くと次のようなものである。
各プロセスはイベントといわれるものを交換し、相互に連携する。アルゴリズムはイベントのハンドラーとして表現される。
ハンドラーは入力となるイベントに対して何かしらの処理を行い、また新しいイベントを発生させることもある。
こういった非同期なイベントのやり取りとして、分散アルゴリズムを表現。

### 用語・記法の整理
* 分散システムは、N 個の異なるプロセスからなる($\prod := \\\{p, q, r, ...\\\}$)
    * 仮定: 各プロセスは互いにお互いを知っている
    * 仮定: この集合は静的で、変化しない
    * 各プロセスは、すべて同じアルゴリズムで動く場合が多い（恐らく、分散アルゴリズムに依るはずだが） → それらの総体として分散アルゴリズムが構成

* すべてのプロセスは1つ以上の modules(or components) から構成される
    * プロセスという抽象化についての詳細は ch2 参照だが、ひとまず計算の単位と思っておけば OK

* module(component) は名前と性質で規定され、入出力となる event がインターフェースとして定義されている
    * 入力イベント: request
    * 出力イベント: indication
        * response のようなものと思って良い

* component が集まってソフトウェアレイヤーを構成する
    * ある component がソフトウェアの特定レイヤーに該当
    * ex. ネットワークやアプリケーションのレイヤーがある。前者は最下層で、後者が最上位層となり、この間に分散アルゴリズムのレイヤーが存在する
    * 同じレイヤーのコンポーネントは event 経由で連携
    * (絵が必要かと)

* event はメッセージとか何でも
    * 表記は < EventType | Attributes,...> あるいは、< co, EventType | Attributes,...>
    * co は component で、これは同じ名前の event が複数の module(component) で利用されるので何の component を考えているかを明記するため

* 各プロセスによる event の処理は、その event の処理を担う component によって行われる
    * こういった component のことを、handler という

event handler の記法パターン
* 記法1: 外部からの event driven なケース
```
upon event < co_1, event_1 | attr_1_1, attr_1_2,...> do
    do something;
    // event_2 を送る
    trigger <co_2, event_2 | attr_2_1, attr_2_2,...>;
```

* 記法2: 内部条件を満たす場合の動作
    * component が内部変数なりをもっていて、それがある条件を満たす時に
    * 内部的なイベントというイメージ
```
upon condition do
    do something;
```

* 記法3: 条件付きの外部 event handler
```
upon event < co, event | attr_1_1, attr_1_2,...> such that conditon do
    do something;
```

ここまでを踏まえ、各プロセスを有限オートマトンと見立てることができ、その集合として分散アルゴリズムが構成される

* （補足）有限オートマトンは5つ組$M := (Q, \Sigma, \delta, q_0, F)$で以下の性質を満たす $M$ として定義
    * $Q$
        * 状態集合
        * finite set
    * $\Sigma$
        * 文字集合
        * finite set
    * $\delta : Q \times \Sigma \rightarrow Q$
        * 遷移関数
    * $q_0 \in Q$
        * 初期状態
    * $F \subset Q$
        * 受理集合
        * 分散システムでプロセスを考える場合は空と思って良い

次にこの分散アルゴリズムが満たすべき重要な性質について説明をする

## 分散アルゴリズムの重要な性質: safety and liveness
* ここの議論は informal なので注意
    * また、良い・悪いの定義も厳密に決まるような話ではないっぽい

* $S$
    * set of process states
* $S^\omega$
    * set of infinite sequences of process states
    * $\sigma \in S^\omega$ を execution という
* $S^*$
    * set of finite sequences of process states
    * $\sigma \in S^*$ を partial execution という
* $\sigma_i$
    * first i states in execution $\sigma$
* $P$
    * subset of $S^\omega$
    * property
    * $\sigma \in P$の時、property を満たすという

この時、以下で propery $P$ に対して safety/liveness を定義

* safety
    * 何も悪いことがおきないという性質
        * ex. deadlock することがない
    * $\forall \sigma \in S^\omega, \sigma \notin P \rightarrow \exists i \in N, \forall \beta \in S^\omega, \sigma_i \beta \notin P$
* liveness
    * 最終的に何か良いことが起きるという性質
        * ex. at least once のメッセージ送信
    * $\forall \sigma \in S^*, \exists \beta \in S^\omega, \sigma \beta \in P$
* 実際に考えるときは別で考えるほうが便利らしいが、アルゴリズムの性質としてはどちらも期待されるものである
    * 例えば、単に safety を満たすというのであればずっと何もしないことも safety である
    * 組み合せることに意味がある

## Reference
* [DEFINING LIVENESS](https://www.cs.cornell.edu/fbs/publications/DefLiveness.pdf)