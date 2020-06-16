---
title: SRE
---

## Service level objectives(SLO)
* Service level objectives(SLO) は、サービス信頼性のターゲットを示す指標
    * 信頼性に関する、データドリブンな意思決定を支えるキーであり ＝ SRE のコアである
    * SLO は、何のエンジニアタスクを優先するかを決定するのを助けるツール
        * ex. For example, consider the engineering tradeoffs for two reliability projects: automating rollbacks
        and moving to a replicated data store. By calculating the estimated impact on our error budget,
        we can determine which project is most beneficial to our users
    * SLO だけ決めて、というだけではなく意思決定に使われることこそが肝(= error budget-based approach)

* 前提
    * エンジニアは巨大企業においてさえ、希少なリソースである
        * ゆえに、重要タスクにこそ時間を投下すべき
    * 機能開発による新規ユーザー獲得と、信頼性やスケーラビリティ向上による既存顧客の幸福というトレードオフ
    * SLO により、信頼性向上タスクの機会費用

* ```without SLOs, there is no need for SREs```

* [可用性とどう向き合うべきか、それが問題だ : CRE が現場で学んだこと](https://cloud.google.com/blog/ja/products/gcp/available-or-not-that-is-the-question-cre-life-lessons)
* [SLO、SLI、SLA について考える : CRE が現場で学んだこと](https://cloud.google.com/blog/ja/products/gcp/availability-part-deux-cre-life-lessons)
* [優れた SLO を策定するには : CRE が現場で学んだこと](https://cloud.google.com/blog/ja/products/gcp/building-good-slos-cre-life-lessons)
    * SLO は四半期ごとにレポートし、四半期の統計に応じてポリシーを決めましょう
    * エラー バジェットを使用した際のページャーへのアラートやリリースの凍結など、SLO に依存するポリシーについては、時間枠を四半期より短くしてください
    * 最後に、特に初期の段階では SLO を共有する相手を意識するようにしてください。SLO は、サービスの期待値を伝えるうえで非常に便利なツールですが、いろんな人に知られてしまうと変更することが難しくなります

### まず目指すべきところ
* error budget-based approach to reliability を実現するために
    * サービスに関して、ステークホルダーが同意した SLO が存在
    * SLO を担保することに責任を持つ人の間で、通常の状況であれば SLO を守れると同意している
    * 組織として、意思決定・優先度付けにエラーバジェットを用いることにコミットしており、エラーバジェットポリシーが承認されている
    * SLO を改良するプロセスがある
* こういった条件を満たさなければ、SLO は単なる KPI やレポーティングのメトリクスになるだけであり、決して意思決定のツールとしては用いられないだろう

### 信頼性のターゲット
* SLO は、サービスのユーザーに提供する信頼性レベルのターゲットである
    * これ以上であれば、ほぼすべてのユーザーは幸福であり、逆に下回る場合は不平やサービスの利用停止を招くかもしれない
* だが、100%目標は誤った目標である
    * SLO をユーザーの満足度と合わせて、100%目標が無理なことは明白
    * （ありえないが）もしサービスが100%信頼性を達成しても、実際にユーザーが感じるのは100%にはならない。なぜなら、サービスまでの何かしらは100%信頼性を担保できないはずだから。
    * 身動きがとれなくなる。新機能とか当然入れると危険しか無いってことになるし。
* ポイント: 兎に角100%はまず無理な目標だということを組織的に握る必要がある。そこから初めて、では何を計測してどんな数値目標を置くのか、という議論にうつれる

### 何を計測するか: SLI
* SLI is an indicator of the level of service
* SLI としてオススメは、2つの数の比として定義するもの。つまり、イベント総数で良いイベントの数を割るというもの。
    * ex. Number of successful HTTP requests / total HTTP requests (success rate)
    * ex. Number of search results that used the entire corpus / total number of search results, including those that degraded gracefully
* なぜこの形式がオススメなのか
    * 0%~100%のレンジに収まり、直感的な解釈が容易
    * the SLO is a target percentage and the error budget is 100% minus the SLO
    * 形式を統一することで、その後の計算とかアラートで使うとか、色々と便利に
        * 要は、分子・分母・閾値を、ネクストの入力として使うことになる
* 更に、SLI specification と SLI implementation に分けて考えることを勧める
    * SLI specification: ユーザーにとって重要となる評価(計測は別で考える)
        * ex. Ratio of home page requests that loaded in < 100 ms
    * SLI implementation: SLI specification とその計測法
        * ex. Ratio of home page requests that loaded in < 100 ms, as measured from the Latency column of the server log. This measurement will miss requests that fail to reach the backend.
        * ex. Ratio of home page requests that loaded in < 100 ms, as measured by probers that execute JavaScript in a browser running in a virtual machine. This measurement will catch errors when requests cannot reach our network, but may miss issues that affect only a subset of users.
* SLI specification と SLI implementation は、1対多の関係で、implementation にはそれぞれメリデメがある
* ポイント: 最初から完璧な SLI/SLO である必要はなく、まずは何かしら設定し、計測を始めることが肝要。そして改善をしていく。
    * ちなみに、SRE 本では現在のパフォーマンスから SLO を決めることは無駄に厳しい SLO になる可能性があるため反対をしたが、
    他に情報がない状況ではスタートとしては悪くはない
* 最初の SLO 群としては、サービスにとってキーとなるような SLI specification から選ぶべし
    * なので、（今更だけど）SLO は1つとは限らない
    * その場合のエラーバジェットってどういう決まり？（どれか1つでも違反したら、だとは思うが）
* ex. SLI をどう決めるか悩む場合には、以下のようにシンプルに始めてみよ。
    1. SLO を定義したいアプリを1つ選ぶ。システムが複数アプリから構成される場合には後で追加していけば良い。
    2. ユーザーを明確に定めよ。誰の幸福を最大化したいのか。
    3. ユーザーがシステムと関わる一般的なやり方を決める。つまり、日常タスクと重要度の高い行動を。
    4. ハイレベルでシステムのアーキテクチャを書く。キーとなるコンポーネント・リクエストやデータのフロー・依存関係を示す。それらコンポーネントを
    纏めて、カテゴリーに分類すべし（カテゴリーについては後述）
* 繰り返しになるが、まずはシンプル（サービスに関係があり、計測も容易）に始めて、磨いていくのが良い

#### コンポーネントの種類
* SLI を決めるのには、システムを、構成するいくつかの種類のコンポーネントに抽象化することが最も容易い
* コンポーネントの種類
    * Request-driven
        * ユーザーはある種のイベントを実行したら、レスポンスが返ることを期待するようなもの
    * Pipeline
        * レコードを入力として、それを変形し、どこかに出力するようなもの
        * リアルタイムの処理もあれば、数時間かかるようなバッチもある
        * ex. ログファイルを読み込んでレポート作成、データを集めて時系列データとアラートを生成するモニタリングシステム
    * Storage
        * データを貯めて利用可能にするもの
* 例は、[A Worked Example](https://landing.google.com/sre/workbook/chapters/implementing-slos/#a-worked-example) を参照
* 種類ごとにどういった SLI が定義されるか、というのは被る部分もある
    * いずれにせよ、サービスのユーザーにとって最も重要と思える要素を5つ程度まで選ぶ
* 典型的なユーザー経験と、ロングテールを同時に考慮したい場合、複数観点評価の SLO もあり
    * ex. 90%のリクエストは < 100ms だが、残りが10s かかるのはやべーシステム
    * そこで、latency SLO として、90% は < 100ms かつ 99% < 400ms という SLO
    * ユーザーの不幸度を測るようなパラメータがある SLI 一般に使えるテクニック
* [コンポーネントの種類と SLI]( https://landing.google.com/sre/workbook/chapters/implementing-slos/#slis-for-different-types-of-services)

### SLI specification から implementation
* 最初はミニマムな作業量で対応可能なところからスタートすべし

### 適切な時間枠の決定
* rolling window と calendar window(e.g., a month)
* rolling window
    * ユーザー経験により関係が強い
    * ex. 月末に大きな動きがあった場合に、そのことをすぐ忘れるユーザーはいないという話
        * cf. calendar window
    * tips: 同じ数の週末を含めるために、整数値の週をオススメ
        * ex. 30日 rolling window だと、週末を4つ含むケースと、5つ含むケースがありうる。で、週末のトラフィックなどが平日と
        異なる場合において、SLI が嬉しくない変化をしてしまうから。
* calendar window
    * ビジネスやプロジェクトにより関係が強い
        * ex. クォーターの SLO 評価から、来 Q の計画を決める
* 期間が短いと意思決定を素早く行える
* 期間が長いと戦略的な意思決定が可能になる
* **general tips**: 4 week rolling window が良い。これにプラス、タスクの優先度付けのために週次のサマリと、プロジェクトプランニングのために
クォーターのサマリを添える

### stakeholder との同意
* 重要
* これが決まれば後はやるだけ。SLO を守るためにモニタリング・アラートを設定しよう

### エラーバジェットポリシーの決定
* SLO が決まると、そこからエラーバジェットを決定できる
* エラーバジェットを利用するには、それを使い果たした時に何をすべきかを述べるポリシーが必要である
    * ポリシーも関係者全員が同意している必要あり
    * 議論があるのであれば、SLI/SLO 設定のレベルから会話していく
* ポリシーを遵守するには、文章化することが必要で、誰が何をするかも書く必要がある

### SLO/エラーバジェットポリシーの文章化
* 関係者が確認できる場所に置かなければならない
* SLO まとめに必要な項目
    * SLO の author/reviewers/approvers
    * date(approved and should next be reviewed)
    * サービスの短い説明
    * SLO の詳細: 目標や SLI implementations について
    * エラーバジェット計算の詳細説明
    * 設定数値の背景
* SLO 文章のレビューは成熟度に応じる
    * 不慣れなうちは頻繁に見直しが必要
* エラーバジェットポリシーまとめに必要な項目
    * policy の author/reviewers/approvers
    * date(approved and should next be reviewed)
    * サービスの短い説明
    * バジェットの使い果たした時に取るべきアクション
    * A clear escalation path to follow if there is disagreement on the calculation or whether the agreed-upon actions are appropriate in the circumstances
    * (想定読者次第だが、エラーバジェット自体の説明もあっても可)
* [Example SLO Document](https://landing.google.com/sre/workbook/chapters/slo-document/)
* [SLO 違反への対処 : CRE が現場で学んだこと](https://cloud.google.com/blog/ja/products/gcp/consequences-of-slo-violations-cre-life-lessons)
    * 参考

### dashboard/reports
* スナップショットの情報として有用
* SLI のトレンドをダッシュボードに出力するのも良い
    * エラーバジェットの消費傾向が通常より高いことがわかったりする
* エラーバジェットが、定量的な議論を可能に

### SLO の継続的な改善
* SLO の改善の「前に」、サービスに対するユーザーの満足度の情報が必要になる
    * ポストやチケット・サポートの電話数
    * ソーシャルメディアの利用
    * サーベイ... などなど
* tips: まずは安く初めて改善していく
* サービスにインシデントなどの問題があった日の何かしら指標（チケット数とか）と SLO(Budget loss) の相関などを分析する
    * 例えば明らかにサービスに問題があったのに SLO に反映できていない -> SLO の coverage に課題がある
    * ちなみに、こういったことは織り込み済み(だから、漸進的な改善が必要)
* こういった場合に取れる手段として
    * SLO の変更
        * SLI が問題を示しているが、SLO による通知などがなされていないなら、SLO を厳しくする必要があるかもしれない
        * SLO の変更で、too many false positives or false negatives が発生するなら、SLI implementation の変更が必要かもしれない
    * SLI implementation の変更
        * よりユーザーに近いところからデータを取得してメトリクスの質を向上
        * coverage を向上させる
* aspirational SLO を考える意味もあり

### SLO とエラーバジェットによる意思決定
* エラーバジェットを使い果たした場合: エラーバジェットポリシーに従って対応
* インシデントの規模の決定に、どれだけのエラーバジェットを消費したか、が使える
    * そうすると、問題の比較に使う事ができ、どの問題に優先的に取り組むかの意思決定に使える！
* あるサービス内での問題比較だけではなく、サービス間で何を優先すべきか、にも使える
    * cf. https://landing.google.com/sre/workbook/chapters/implementing-slos/
    * SLO decision matrix

### Advanced
* Modeling User Journeys
* Grading Interaction Importance
* Modeling Dependencies
* Experimenting with Relaxing Your SLOs

## 実例: VALET
* [SLO Engineering Case Studies](https://landing.google.com/sre/workbook/chapters/slo-engineering-case-studies/)
* Automating VALET Data Collection を読む
    * Note that our SLOs are a trending tool that we can use for error budgets,
    but aren’t directly connected to our monitoring systems.
    Instead, we have a variety of disparate monitoring platforms, each with its own alerting.
    Those monitoring systems aggregate their SLOs on a daily basis and publish to the VALET service for trending.
    The downside of this setup is that alerting thresholds set in the monitoring systems aren’t integrated with SLOs; however, we have the flexibility to change out monitoring systems as needed.

## SLO とアラート
### アラートで考えること
* SLO における重要イベント: error budget の多くを消費するイベント
* SLO に関するアラートでは、それを通知されることが目的となる
* アラート戦略の評価指標
    * Precision
    * Recall
    * Detection time
        * How long it takes to send notifications in various conditions. Long detection times can negatively impact the error budget.
    * Reset time
        * How long alerts fire after an issue is resolved. Long reset times can lead to confusion or to issues being ignored.
        * これはこのあと見ていけばよく分かる

### 重要イベントをアラートするために
* 繰り返しだが、今回の重要イベントは、エラーバジェットを多く費消してしまうようなものを指す

### Target Error Rate ≥ SLO Threshold
* 以下では、reporting period := 30days
* if the SLO is 99.9% over 30 days, alert if the error rate over the previous 10 minutes is ≥ 0.1%
    * 検知は速いが、precision が低い
* error rate is equal to the SLO のときにアラートということは、エラーバジェットを alerting window size / reporting period 費消したということ
    * ex. 0.1% error rate for 10 minutes だと、10 / (60 * 24 * 30) * 100 = 0.02%
* (1 - SLO) / error ratio * alerting window size = detection time で、alerting window size / detection time = error ratio / (1 - SLO)
    * ex. 100% outage(error ratio=1)だと、(1-0.999) / 1 * 10(m) * 60 = 0.6s
    * error ratio は、alerting window size におけるエラーの割合
    * 144 = 6(10分に1回) * 24(時間)

### Increased Alert Window
* 期間のエラーバジェットのどれだけを消費した時に通知するか、という基準で考える
* 検知は相変わらず速く、precision も改善
    * より長い期間の error rate を考慮するので
* 一方 reset time が非常に悪く、計算負荷も増大
* 1 - SLO を error ratio で消費するのに必要な alerting window size で計算
    * 1 - SLO と error ratio が等しい場合に、(上にも書いたが) 1 - SLO より何倍 ratio があるか、で割れば良い
    * ex. SLO: 99.9% alerting window size: 36h の場合、100% outage(error ratio = 1)だと、
    (0.001 / 1) * 36 = 0.036h で、0.036h * 60 = 2.16m

### Incrementing Alert Duration
* 同じく precision は良い
* recall が悪く、detection time も良くない
    * duration は、インシデントの緊急度で比例してあがるというものではないので
    * ex. a 100% outage alerts after one hour, the same detection time as a 0.2% outage だが、前者はバジェットの140%を食いつぶす
        * 1時間で、100%outrage なので、1 / (30 * 24 * 0.001) = 1.38 ~ 140%程度
            * 最初の1は、一時間基準で。error budget 量で考えて単位は1とすれば分母が1ヶ月での総量(h 単位)で、分子が100%の場合に1時間で消費した量
    * 境界を行ったり来たりするような場合も、検知されず
* spike の計算
    * (1 / 12) / (30 * 24 * 0.001) = 0.115 ~ 12%程度を spike で消費
        * 1/12 = 5分

### Alert on Burn Rate
* detection time と precision の両立のために burn rate という概念を導入
    * burn rate: SLO に対して、エラーバジェットを費消する速さがどれだけなのか
* burn rate は、window size と何%のエラーバジェット費消でアラートにするかを決めることで計算される
    * ex. window size: 1h で5%エラーバジェット費消を検知するとした場合、
    reporting period: 30days なので、b(burn rate) * 1h / 30 * 24(h) = 0.05 -> b = 0.05 * 30 * 24 = 36
    となり、burn rate 36
* (1 - SLO) / error ratio * alerting window size * burn rate がアラート発火までの時間
* burn rate * alerting window size / period が、アラート発火までに費消したエラーバジェット
* 1000 * x(h) / (30 * 24) = 1 となる x = 0.72(h)で、0.72 * 60 = 43m
* Reset time: 58 m
    * 36 * 0.001 * 100 で、3.6%がアラートの基準であった
    * 60m window で、2m の error rate が1(100%)、それ以外0だと考えると、結局58m 経過した時点での
     2 / 60 = 0.033..で、アラートの基準を割る(エラー・リクエストの絶対数がいくらかは関係ない)
* A 35x burn rate never alerts, but consumes all of the 30-day error budget in 20.5 hours.
    * 35 * x / 30 * 24 = 1 -> x = 20.57...
        * burn rate 基準で考えれば、SLO の % は計算には関係しない
* shorter time window で利用できるし、detection time も良い
* 一方で、low recall(ex. 上の35burn rate)・reset time もまだ長い

### Multiple Burn Rate Alerts
* アイディア: Fast burning は素早く検知され、Slow burning はより長い time window を必要とする
    * https://developers.soundcloud.com/blog/alerting-on-slos
* 複数の burn rate/time window で、burn rate の利点を持ちつつ、lower error rate も見逃さないようにする
* start point としてオススメ:
    * paging: 2% budget consumption in one hour and 5% budget consumption in six hours
    * tickets: 10% budget consumption in three days as a good baseline
* b = 0.02 * (30 * 24 * 0.1)
    * b / (30 * 24) = 0.02(=2%) となるような b を求めると、b = 14.399.. つまり burn rate: 14.4となる
* b = 0.05 * (30 * 24) / 6 で b = 6.0 から burn rate: 6
    * b * 6(h) / (30 * 24) = 0.05から
* good precision/recall で、色々調整も効く
* 一方で、調整すべきパラメータが増え、複数アラートが発火する関係で、suppression も必要に。
また、reset time は 3day といった長期間のものをいれた関係で長くなる

### Multiwindow, Multi-Burn-Rate Alerts
* 今実際に budget burning しているときにのみ知らせてほしい
    * short window のチェックを組み込む
* short window の長さは long window の 1/12 の長さがベースライン
* (0.001) / 0.15 * 60 * 14.4 = 5.76m
* parameter recommendation
    * severity/long window/short window/burn rate/error budget consumed
    * Page/1 hour/5 minutes/14.4/2%
    * Page/6 hours/30 minutes/6/5%
    * Ticket/1 days/2 hours/3/10%
    * Ticket/3 days/6 hours/1/10%
* good precision/recall で flexible
* 一方で、複雑
    * request の種別でいくつか group 化けて、同じグループ内での parameter 共通化という方法もあり
    * cf. SRE workbook ch5. Alerting at Scale

### Extreme Availability Goals
* 高い・低い、いずれの場合においても少し特殊な扱いが必要となる
* 高い場合だと、秒でバジェットを食いつぶすこともある
    * canarying release などが必要に

## Runbook/Playbook
* [maintaining-playbooks](https://landing.google.com/sre/workbook/chapters/on-call/#maintaining-playbooks)
    * 本番環境の変化に合わせ、playbook の詳細も古くなってしまう
    * 大きな思想の違い
        * general な内容に留めて変化による影響を減らすことを目指す思想
            * ex. RPC Errors High というアラートがあれば、全てに共通の1つだけ書く
        * step-by-step な playbook で、人差や MTTR を減らすことを目指す思想
    * チーム内でここが統一出来ていないと、playbook 自体がハチャメチャに
    * 決めるのが難しいので、ミニマムの構造は少なくとも決定しておいて、これを超える情報を逐一チェックするというやり方もある
        * playbook に、アラート時にいつもやるコマンドのリスト、的なものが増えてきたら自動化すべし

## Reference
* [awesome-sre](https://github.com/dastergon/awesome-sre)
* [The SRE I aspire to be](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_aknin.pdf)
* [Evolution of Observability Tools at Pinterest](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_abbas.pdf)
* [Hiring Great SREs](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_rutkin.pdf)
* [SRE & Product Management](https://www.usenix.org/sites/default/files/conference/protected-files/srecon19emea_slides_wohlner.pdf)
* [Training Site Reliability Engineers: What Your Organization Needs to Create a Learning Program](https://landing.google.com/sre/resources/practicesandprocesses/training-site-reliability-engineers/)
* [サイト信頼性エンジニアリングのドキュメント](https://docs.microsoft.com/ja-jp/azure/site-reliability-engineering/)
  * リンク集

### Runbook
* [Writing Runbook Documentation When You’re An SRE](https://www.transposit.com/blog/2020.01.30-writing-runbook-documentation-when-youre-an-sre/)
    * Tips and tricks for writing effective runbook documentation when you aren’t a technical writer

### Game day
* [AWS Game Day](https://jon.sprig.gs/blog/post/1238)
* [SRE として Adversarial Game Day (敵性ゲームデイ) を行う方法](https://blog.newrelic.co.jp/best-practices/how-to-run-a-game-day/)