---
title: Monitoring
---

## Monitoring
### 定義や分類
定義はいろいろあるなかで1つの考え方として、次の4つの作業に分けるというものがある
(cf. 入門 Prometheus)
* アラート
* デバッグ情報の提供
* トレンド調査
    * 設計判断やキャパシティプランニングの根拠として
* 部品提供
    * モニタリングシステムの適切な部分を他の目的に転用

モニタリングの対象はイベント
* イベントはほとんどあらゆるもの(草)
* イベントにはコンテキストが含まれ、すべてのコンテキストがあればデバッグにも、技術的・ビジネス的なパフォーマンスの理解にも役立つ
    * ex. HTTP リクエストには、送信元・送信先の IP アドレス・URL・Cookie などのコンテキストがある
* しかしながら、すべてのコンテキストとなると量が莫大なので削減する方法が求められる
    * プロファイリング
        * いつもすべてのイベントのすべてのコンテキストを保持することは出来ないが、限られた時間であれば一部のコンテキストは残せるという考え方をとる方法
        * ex. eBPF
            * http://www.brendangregg.com/ebpf.html
        * 目的は短期的なデバッグにあることが多い
    * トレーシング
        * すべてのイベントではなく、例えば注目している関数を通過する100回のうちのイベント1回など、イベントの一部のみをみる
            * サンプリングでデータ量を合理的な範囲に抑える
    * ロギング
        * イベントの一部を対象として、それらのイベントそれぞれのコンテキストの一部を記録
            * ex. 送られてくる HTTP リクエスト、発行するすべてのデータベース呼び出し
            * リソースを過度に消費しないように、ログエントリあたりのフィールド数は100程度までに制限される
        * (通常)イベントのサンプリングはしないので、フィールド数に制限があるといえ判断に大いに役立つ
        * ロギングも十人十色な定義をされがちだが、1つの分類として
            * トランザクションログ
                * 重要なビジネス記録で、何が何でも残さないといけないもの
            * リクエストログ
                * ex. すべての HTTP リクエストやすべてのデータベース呼び出しの記録など
                * 実装や内部の最適化につかわれることが多く、一部無くなってもそれほど重大ではない
            * アプリケーションログ
                * プロセス自体についてのログ
                * 人間が直接読むことが多いため、通常の運用では1分あたり数エントリ程度に抑えたい
            * デバッグログ
                * 非常に詳細で、作成・格納コストが高くつくことが多い
                * プロファイリングに近く、信頼性や保存の要件は低め
        * これら、信頼性などの非機能的な要件が異なるため、同じように扱うべきではない
            * 特にシステムが大きくなってきた場合、デバッグログは切り分けて他のログとは別に処理できるようにすべき
            * また、ユニットテストの対象とする場合、トランザクションログと一部リクエストログ程度になるであろう
    * メトリクス
        * (基本的に)コンテキストを無視して様々なタイプの集計を経時的に管理
            * こちらも制限として1プロセスあたり1万が妥当な上限
        * 常にコンテキストが無視されるというわけではない
            * ex. HTTP リクエストの場合の URL パス
            * 注意: コンテキスト分だけメトリクスが増えるので、1万という（現実的に妥当とされる）制限に注意
     
* ログとメトリクスの対照
    * メトリクス: プロセス全体のイベントを集められるが、コンテキストのカーディナリティは制限されている
    * ログ: 一つのイベントタイプについてすべての情報を集められるが、カーディナリティに制限のない百個のフィールドを追跡するのみ

### ダッシュボード
* 目的が複数あるようなダッシュボードは避ける
    * 1つに1つ
* グラフ数は抑える

## Prometheus
* 上の分類上は、メトリクスベースのモニタリングシステム
* pull 型

### metrics
* メトリクスが役立つかどうかの基準
    * メトリクス全体の合計か平均を取って意味のある結果が得られるかどうか
        * ex. パス別・メソッド別の HTTP リクエストカウンタの合計はリクエスト総数となり、意味がある
        * 例外. サマリメトリクスの分位数
            * テーブル例外といわれる
            * メトリクスによって集計ができずとも、名前が何らかの正規表現にマッチするメトリクスを抽出するくらいであれば、
            ラベルを使った（濫用した）ほうがましというケース
            * 何らかの正規表現にマッチする名前をもつメトリクスの抽出というようなことは、グラフやアラートでは絶対にしてはいけない

#### 種別
* counter
    * 知りたいのはどのようなペースで増えるか
* gauge
    * 何らかの状態のスナップショットの値
    * 知りたいのは実際の値
* summary
    * イベントの何らかのサイズ(非負数値)
* histogram

#### 命名
* <ライブラリ>\_<名前>\_<単位>\_<サフィックス>
    * 注意: ラベルにするようなものを入れてはいけない
        * PromQL でそのラベルの違いを無視して集計したときに正しくなくなるから
* _total/_count/_sum/_bucket のサフィックスにはそれぞれ意味があるのでメトリクスのサフィックスとしては使わないようにすべき
    * それぞれ、カウンタ・サマリ・ヒストグラムあたりで使う
* 名前の最後にメトリクスの単位をいれることも推奨
    * ex. myapp_requests_processed_bytes_total
    * Prometheus での時間の基本単位は秒

#### 提供
* ダイレクトインストルメンテーション
    * コード中にクライアントライブラリを埋め込み、メトリクス生成装置を追加
* exporter
    * コードをいじれない(アクセスもできないとか)、場合にもメトリクスにアクセスするためのインターフェースを持つことは多い
    * そのようなときに、アプリと一緒にデプロイし、アプリケーションから必要なデータを収集・Prometheus に返すもの

### instrumentation
* 何を対象とすべきかとどの程度の量を計測すべきかを順に書く

#### サービスの instrumentation
* サービスを大きく分類すると
    * オンライン配信システム
        * 人間や他のサービスがレスポンスを待っているもの
            * ex. web server, db
        * 重要メトリクスは、Request(Rate), Errors, Duration で RED
    * オフライン配信システム
        * 待っているものがないもの
            * ex. ステージ間にキューがあるログ処理システム
        * 重要メトリクスは、Utilization(利用状況), Saturation(飽和), Errors で USE
            * Utilization はサービスがどの程度稼働しているか、Saturation はキューの処理待ち量など
        * バッチを使う場合には、バッチのためのメトリクスと個々の要素のメトリクスの両方を用意すると役立つ
    * バッチジョブ
        * オフライン配信システムと似ているが、継続的な実行でなく定期実行という点が異なる
            * つまり常に実行されているわけではない → スクレイピングがあまり役に立たない
                * pushgateway などが使われる
        * バッチジョブ終了時に実行にかかった時間、ジョブの各ステージでかかった時間、ジョブが最後に成功した時刻など

#### ライブラリの instrumentation
* 小規模なサービスと考えて良い
* ライブラリの多くは、オンライン配信システムのサブシステム = 同期関数呼び出しなので、RED が役立つ
* エラー発生箇所やロギング箇所へのメトリクス追加も役立つ
* スレッドやワーカープールも、オフライン配信システムと同様に、計測すべき

#### 量はどう決めるべきか
* 

* 成功と失敗よりも、全体と失敗をとると割合の計算などがすぐ出来て便利

### label
* ラベルの種類として、instrumentation label と target label が存在
    * 前者は instrumentation で指定するもの
        * ex. HTTP リクエストのタイプ、やり取りするデータベースの名前、その他システム内の具体的な情報といった、アプリ・ライブラリ内部でわかっているもの
        * instrumentation label を途中で追加・削除すると、互換性が失われることに注意
    * 後者はモニタリングの特定の対象(Prometheus がスクレイプする対象)を識別するもの
        * ex. 環境(stg/prd)や、どのインスタンスか
        * アプリケーションや exporter ではなく Prometheus で設定
            * チームごとに自分たちに有用なラベル階層をつくることができる
* ラベルはメトリクスと異なり文字列なので、ログベースのモニタリングに近づくような利用法ができる
* ユースケース一覧
    * 列挙
        * ex. ゲージに状態のためのラベルを追加し、メトリクスを0/1のブール値とする
            * avg_over_time(resource_state[1h]) で各状態の時間割合を計算したり
            * sum without(resource)(resource_state) で各状態ごとのリソース数がわかる
    * info metrics(or machine roles)
        * バージョン番号などのビルド情報のアノテーションとして有用
            * ターゲットのアノテーションとして使いたいすべての文字列をラベルとし、値を1とするゲージを使う
        * _info というサフィックスを付ける慣習

* instance/job はターゲットが必ずもつラベル

#### 注意
* メトリクスのラベル名はアプリケーションプロセスの生涯を通じて変更してはいけない
    * 変更の必要性を感じる場合には、そのユースケースではログベースのモニタリングシステムが適切かもしれない
    * ここでいうラベルの種類は？
* ラベルのつけすぎには注意
    * 単純にコストがかかる
    * Prometheus 1.x では、およそ1千万くらいの時系列数が性能限界の目安
        * 時系列数が重要なので、例えばラベルが多そうだから、メトリクスに適当な名前をつけて別扱いしよう、といったことをしても、使いにくくなるだけで何も意味がない
        * 例えばそこから、メールアドレス・顧客・IP アドレスなどはラベルとして不適切
    * 厳しそうであれば、ログベースのモニタリングに切り替える時かもしれない

#### 命名
* instrumentation label を定義するときに、env, cluster, service, team, zone, region のような target label で使われそうなものは避ける

### 長期記憶
* Remote read/write API + 他のシステムとの連携
* API に対応している長期記憶ストレージ + InfluxDB, Cortex, Thanos, M3DB

### scrape
* スクレイピングの失敗はアプリケーションログには書き込まない
    * スクレイプエラーは、Targets ページ or デバッグログ有効化(--log.level debug をコマンドラインで指定)
* 一般に、 Prometheus では健全性を保つために単一のスクレイプインターバルを使うようにすべき

## Grafana

### annotations

### Dashboard as code
* [grafana-dashboards](https://github.com/adamwg/grafana-dashboards)

## モニタリングの意義
* Alert on conditions that require attention.
* Investigate and diagnose those issues.
* Display information about the system visually.
* Gain insight into trends in resource usage or service health for long-term planning.
* Compare the behavior of the system before and after a change, or between two groups in an experiment.

## 症状と原因(Symptoms vs Causes)
* モニタリングシステムが答える必要のある2つの疑問
    * 何が壊れたのか・なぜそれが壊れたのか
* 「何が」と「なぜ」の区別はモニタリングルールを書く上で最も重要な区別
* ちなみに昔は、この区別はあまり意識することはなかった
    * ex. web server down(cause) で、サイトダウン(symptom)なので、単にアラート
    * 今は 1 server down でページレベルのアラートを出すことはないこともごく一般的な話
* Paging alerts, those that wake you up in the night, should be based on symptoms, i.e. something that actively or imminently impacts the user experience
* この辺りは、Reference の My Philosophy on Alerting が詳しい

## ブラックボックスとホワイトボックス
* Google では、ホワイトボックスモニタリングを多用し、一部重要な部分でブラックスボックスモニタリングを行っている
* ブラックボックスは「現在、システムが正常に動作していない」というような症状を扱う
    * symptom-oriented なので、ページングとして人の対応が必要というのを強制できるメリットがある
* ホワイトボックスは、ログや HTTP エンドポイントなど、システム内部を調査する機能に依存
    * よって、リトライでマスクされてしまっているような障害や、近々生じそうな問題の検出も可能
* 多層的なシステムに於いては、ある人にとっての symptom が別の人にとって、 cause であることがある
    * ex. 低速な DB 読み込みは、DB SRE にとっては symptom で、frontend SRE にとっては、cause
    * よって、ホワイトボックスモニタリングは symptom-oriented なときもあり、cause-oriented なときもある

## USE method

## 【事例】gitlab
* [runbooks](https://gitlab.com/gitlab-com/runbooks)

## Reference
* [My Philosophy on Alerting](https://docs.google.com/document/d/199PqyG3UsyXlwieHaqbGiWVa8eMWi8zzAn0YfcApr8Q/edit)
    * 考え方がまとまっていて良い
    * Symptom/Cause の話もそうだが、 Tickets, Reports and Email/ Tracking & Accountability とかも必読
* [The RED Method](https://grafana.com/files/grafanacon_eu_2018/Tom_Wilkie_GrafanaCon_EU_2018.pdf)
* [The USE Method](http://www.brendangregg.com/usemethod.html)
* [what-makes-a-good-runbook](https://www.transposit.com/blog/2019.11.14-what-makes-a-good-runbook/)
* [The SRE I aspire to be](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_aknin.pdf)
* [](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_abbas.pdf)
* [Hiring Great SREs](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_rutkin.pdf)
* [SRE & Product Management](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_wohlner.pdf)
* [The Factors That Impact Availability, Visualized](https://www.vividcortex.com/blog/the-factors-that-impact-availability-visualized)
  * Availability と MTBF/MTTR の関連について
* [Monitoring Theory](http://widgetsandshit.com/teddziuba/2011/03/monitoring-theory.html)
  * informative/actionable の定義・分類と実例
  * とても良い
* [Customizing alertmanager notifications](https://promcon.io/2018-munich/slides/lightning-talks-day2-01_customizing-alertmanager-notifications.pdf)
  * 本題は兎も角 ticket 使えという話
* [The Art of SLOs](https://landing.google.com/sre/resources/practicesandprocesses/art-of-slos/)
  * SLO の Workshop 資料
* [モダンなシステムにSLI/SLOを設定するときのベストプラクティス](https://blog.newrelic.co.jp/engineering/best-practices-for-setting-slos-and-slis-for-modern-complex-systems/)
  * 良い
* [Datadog を利用して SLO を管理しよう！](https://dev.classmethod.jp/articles/slo-monitoring-tracking-by-datadog/)
  * 参考として

### 実装例
* [dispatch](https://github.com/Netflix/dispatch)
  * All of the ad-hoc things you're doing to manage incidents today, done for you, and much more!
  * Netflix のインシデント管理ツール
