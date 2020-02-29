---
title: Prometheus
---

## Prometheus

## metrics
* [How does a Prometheus Counter work?](https://www.robustperception.io/how-does-a-prometheus-counter-work)
* [How does a Prometheus Gauge work?](https://www.robustperception.io/how-does-a-prometheus-gauge-work)
* [How does a Prometheus Summary work?](https://www.robustperception.io/how-does-a-prometheus-summary-work)

## functions
* hash: 8fdfa8abeaaaf021d3ca614e4b50af317b1f4145
* https://github.com/prometheus/prometheus/blob/master/promql/functions.go

``` go
// extrapolatedRate is a utility function for rate/increase/delta.
// It calculates the rate (allowing for counter resets if isCounter is true),
// extrapolates if the first/last sample is close to the boundary, and returns
// the result as either per-second (if isRate is true) or overall.
func extrapolatedRate(vals []Value, args Expressions, enh *EvalNodeHelper, isCounter bool, isRate bool) Vector {
	ms := args[0].(*MatrixSelector)

	var (
		matrix     = vals[0].(Matrix)
		rangeStart = enh.ts - durationMilliseconds(ms.Range+ms.Offset)
		rangeEnd   = enh.ts - durationMilliseconds(ms.Offset)
	)

	for _, samples := range matrix {
		// No sense in trying to compute a rate without at least two points. Drop
		// this Vector element.
		if len(samples.Points) < 2 {
			continue
		}
		var (
			counterCorrection float64
			lastValue         float64
		)
		for _, sample := range samples.Points {
		    // comment: counter で0 に戻ったときの対策として、counterCorrection をおいているのだと思われ 
			if isCounter && sample.V < lastValue {
				counterCorrection += lastValue
			}
			lastValue = sample.V
		}
		resultValue := lastValue - samples.Points[0].V + counterCorrection

		// Duration between first/last samples and boundary of range.
		durationToStart := float64(samples.Points[0].T-rangeStart) / 1000
		durationToEnd := float64(rangeEnd-samples.Points[len(samples.Points)-1].T) / 1000

		sampledInterval := float64(samples.Points[len(samples.Points)-1].T-samples.Points[0].T) / 1000
		averageDurationBetweenSamples := sampledInterval / float64(len(samples.Points)-1)

		if isCounter && resultValue > 0 && samples.Points[0].V >= 0 {
			// Counters cannot be negative. If we have any slope at
			// all (i.e. resultValue went up), we can extrapolate
			// the zero point of the counter. If the duration to the
			// zero point is shorter than the durationToStart, we
			// take the zero point as the start of the series,
			// thereby avoiding extrapolation to negative counter
			// values.
			// comment:  
			durationToZero := sampledInterval * (samples.Points[0].V / resultValue)
			if durationToZero < durationToStart {
				durationToStart = durationToZero
			}
		}

		// If the first/last samples are close to the boundaries of the range,
		// extrapolate the result. This is as we expect that another sample
		// will exist given the spacing between samples we've seen thus far,
		// with an allowance for noise.
		extrapolationThreshold := averageDurationBetweenSamples * 1.1
		extrapolateToInterval := sampledInterval

		if durationToStart < extrapolationThreshold {
			extrapolateToInterval += durationToStart
		} else {
			extrapolateToInterval += averageDurationBetweenSamples / 2
		}
		if durationToEnd < extrapolationThreshold {
			extrapolateToInterval += durationToEnd
		} else {
			extrapolateToInterval += averageDurationBetweenSamples / 2
		}
		resultValue = resultValue * (extrapolateToInterval / sampledInterval)
		if isRate {
			resultValue = resultValue / ms.Range.Seconds()
		}

		enh.out = append(enh.out, Sample{
			Point: Point{V: resultValue},
		})
	}
	return enh.out
}

// === rate(node ValueTypeMatrix) Vector ===
func funcRate(vals []Value, args Expressions, enh *EvalNodeHelper) Vector {
	return extrapolatedRate(vals, args, enh, true, true)
}
```

## PromQL
* or
    * これの結果は、条件一致する最も高いものが、alert 毎に返る
```
ALERTS{environment="develop", severity="critical", alertstate="firing"} * 10000
or
ALERTS{environment="develop", severity="warning", alertstate="firing"} * 100
or
ALERTS{environment="develop", severity="notice", alertstate="firing"}
```

* =~
    * {} の中に label の条件を入れると、そのラベルを持つものだけが返される

* [Query variable](https://grafana.com/docs/grafana/latest/features/datasources/prometheus/#query-variable)
* ex. 1連の例
    * variables: job が label_values(node_uname_info, job)
        * label_values(metric, label): Returns a list of label values for the label in the specified metric.
        * これで、job の label の値一覧
    * variables: name が label_values(node_uname_info{job=~"$job"}, nodename)
        * 同じく、nodename の値一覧だが、上記の job に依存
    * variables: node が label_values(node_uname_info{nodename="$name"}, instance) + Regex: /([^:]+):.*/
        * 同じく、instance の値一覧で、上記の name に依存
            * ex. node-exporter:9100 が instance で、Regex から、node-exporter 部分
    * variables: port が label_values(node_uname_info{instance=~"$node:(.*)"}, instance) + /[^:]+:(.*)/
        * 上記の例だと 9100 の port 部分

## node exporter
* node_disk_write_time_seconds_total{ictsrv="devops"}[5m] の結果は例えばこんな感じ
    * つまり、スクレイプするたびに、write にかかった時間の秒数総和が得られる
        * 恐らくこのデータだと、毎回スクレイプでも write が完了してないので、value が変わっていないのかと

## blackbox exporter
* up(blackbox exporter 自体)と、probe_success が両方 alert として必要なはず
    * それぞれ別で
* [Checking if SSH is responding with Prometheus](https://www.robustperception.io/checking-if-ssh-is-responding-with-prometheus)

## recording rules sample

### node exporter
* instance:node_cpus:count
    * count without (cpu, mode) (node_cpu_seconds_total{mode="idle"})
        * node_cpu_seconds_total{mode="idle"}: これで、idle 状態の秒数が CPU コア毎にインスタンスベクトルで
        * count without (cpu, mode) (node_cpu_seconds_total{mode="idle"}): これで、cpu/mode を無視してカウント
    * The count of CPUs per node, useful for getting CPU time as a percent of total.

* instance_cpu:node_cpu_seconds_not_idle:rate1m
    * sum without (mode) (rate(node_cpu_seconds_total{mode!="idle"}[1m]))
        * rate(node_cpu_seconds_total{mode!="idle"}[1m]): idle 状態以外の rate
        * sum without (mode) (rate(node_cpu_seconds_total{mode!="idle"}[1m])): mode を無視して、core 毎の sum
    * CPU in use by CPU
        * mode は纏められる

* instance_mode:node_cpu_seconds:rate1m
    * sum without (cpu) (rate(node_cpu_seconds_total[1m]))
        * rate(node_cpu_seconds_total[1m]): 全ての状態で rate
        * sum without (cpu) (rate(node_cpu_seconds_total[1m])): cpu core を無視して sum
    * CPU in use by mode
        * core は纏められる

* instance:node_cpu_utilization:ratio
    * sum without (mode) (instance_mode:node_cpu_seconds:rate1m{mode!="idle"}) / instance:node_cpus:count
        * sum without (mode) (sum without (cpu) (rate(node_cpu_seconds_total{mode!="idle"}[1m]))) / count without (cpu, mode) (node_cpu_seconds_total{mode="idle"})
            * これに等しいはず..({mode!="idle"} の位置に注意 [1m]を先に使うと、結果が〜なので、selector は使えないとかのはず)
        * sum without (mode) (instance_mode:node_cpu_seconds:rate1m{mode!="idle"}): 分子で、idle mode 以外の sum
        * 分母は、CPU core 数
    * CPU in use ratio
        * idle を除いた CPU 使用率(core を纏めて)

* job:node_cpu_utilization:avg_ratio
    * avg without (fqdn, instance) (instance:node_cpu_in_use:ratio)
        * avg without (fqdn, instance) (sum without (mode) (sum without (cpu) (rate(node_cpu_seconds_total{mode!="idle"}[1m]))) / count without (cpu, mode) (node_cpu_seconds_total{mode="idle"}))
            * これに等しい
        * fqdn/instance は、例えば hoge.com に対して、インスタンスが複数というケースが想定されて、それを纏めて無視(グループ化)
            * ただそうなると、何故 fqdn だけでは駄目かというと？
                * 1 instance が複数 fqdn に関与している可能性も考えていたり？ わからん..
    * CPU summaries

* instance:node_memory_available:ratio
    * (node_memory_MemAvailable_bytes or (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes))
    / node_memory_MemTotal_bytes
        * 入門 Prometheus を読めばわかるが、node_memory_MemAvailable_bytes がカーネルが推計した本当に使えるメモリ容量らしい
        * が、追加が Linux 3.14 らしいので、それ以外の場合を考慮するために or 後ろがあるはず
            * これらの説明も本に書いてある
        * node_memory_MemTotal_bytes はマシンの物理メモリ総量
    * 利用可能メモリ割合

* instance:node_memory_utilization:ratio
    * 1 - instance:node_memory_available:ratio

* instance:node_filesystem_avail:ratio
    * node_filesystem_avail_bytes{device=~"(/dev/.+|tank/dataset)"} / node_filesystem_size_bytes{device=~"(/dev/.+|tank/dataset)"}
        * tank/dataset とかはそれぞれの環境によりそう
    * node_filesystem_avail_bytes{device=~"(/dev/.+)"} / node_filesystem_size_bytes{device=~"(/dev/.+)"}
    * 利用可能ディスクスペースの割合

* instance:node_disk_writes_completed:irate1m
    * sum(irate(node_disk_writes_completed_total{device=~"sd.*"}[1m])) without (device)
        * やっぱりこれの device も環境によるのでは？
    * sum(irate(node_disk_writes_completed_total{device=~".*"}[1m])) without (device)
    * device 関係なく、完了した書き込み I/O の数？

* instance:node_disk_reads_completed:irate1m
    * sum(irate(node_disk_reads_completed_total{device=~"sd.*"}[1m])) without (device)
    * read ver

## Label 設計
* ラベルの生成流れ
### ターゲット作成まで
* [Life of a Label](https://www.robustperception.io/life-of-a-label)
1. service discovery
2. relabel_configs を反映・drop/keep actions も
3. __address__ の port 判定や、__meta_ ラベルの除去
4. instance ラベルの設定
5. ターゲット作成

### スクレイピング処理
1. URL を生成
    * scheme: __scheme__ host: __address__ path: __metrics_path__ params: __param_* + config
2. スクレイプ実行
3. スクレイプから得られる個別のラベル(インストルメンテーションラベル)について
    * ターゲットラベルに存在しない → 追加
    * ターゲットラベルに存在
        * → honor_labels が設定されている → インストルメンテーションラベルを優先し、追加
        * → honor_labels が設定されていない → exported_ という prefix を該当ラベルに追加
4. metric_relabel_configs を適用
5. up と scrape_duration_seconds with target labels で終了

### ターゲットラベル
* ターゲットが何かを教えてくれるラベル
* service discovery + relabeling で作られる
    * 特に後者を使うことで、組織にとって意味のあるラベルを作る
* ターゲットラベル = identity of the target なので、変更すると大変
    * 勿論、収集対象システムの変化に応じて変わりうるが、戦略的にやる必要あり
* ラベル数は minimal に
    * 増やすほど、グラフ・アラートでの利用で意識することが増える
* ターゲットの identity であることを忘れず、ラベルをつかうことで least to least to most specific
へ、各ラベルによって他のターゲットと、そのラベルがないと出来ないような方法で区別出来るように
    * 何をモニタリングしている・何のスタックなのかという情報であれば、大抵は job ラベルで
    * instance ラベルがあれば、一般にターゲットを一意に特定可能なはずである
        * なので追加で、host/alias は不要なはず
        * instance により命名をというのはアリ
    * よって2つ合わせて the instance label uniquely identifies a target within a job となる

* [Controlling the instance label](https://www.robustperception.io/controlling-the-instance-label)
    * instance ラベルについて
    * ec2 だと、Name が候補だけどこれだと一意には無理な場合がありそうなので、
        * デフォルトの __address__ つまり、private_ip:port で妥協
        * instance_id
        * Name + 一意な何か
            * 解釈もしやすくて良いかも
            * instance = .. でクエリを手打ちすることはほぼない気がするので、少々長くても許される？
* [Target labels are for life, not just for Christmas](https://www.robustperception.io/target-labels-are-for-life-not-just-for-christmas)
    * target label はモニタリング期間で constant である必要
    * software version とかの、モニタリング期間中にターゲットに関して変化するものでグルーピングしたい場合はどうすれば？
        * [How to have labels for machine roles](https://www.robustperception.io/how-to-have-labels-for-machine-roles)
        * [Exposing the software version to Prometheus](https://www.robustperception.io/exposing-the-software-version-to-prometheus)
* [why-cant-i-use-the-nodename-of-a-machine-as-the-instance-label](https://www.robustperception.io/why-cant-i-use-the-nodename-of-a-machine-as-the-instance-label)
    * シンプルに、ターゲットラベルは service discovery + relabeling 段階で決まるもので、つまりスクレイピングが実際に行われる前だから
* [Target labels, not metric name prefixes](https://www.robustperception.io/target-labels-not-metric-name-prefixes)
    * Services are not distinguished by their metric names in Prometheus.
    * 一般に Prometheus ではメトリクスの名前はアプリやデプロイには紐付かない
        * これによって、様々なアプリに渡って集計が出来る
        * ex. RPC library の作者からは、そういった様々なアプリの集約が出来ると非常に嬉しい
        * メトリクスを使うのは、アプリの開発者だけではない！
    * では、例えば process_cpu_seconds_total をどうやって自分の気にするアプリのものからだと判断するのか？
        * ターゲットラベルを使う！
            * job/instance label
    * So if you're trying to prefix all the names on a /metrics, look to target labels instead.

### Job Label 設計
* job ラベルは、同じ目的を持つ一連のインスタンスを示し、一般にすべて同じバイナリと構成で実行される。instance ラベルは、ジョブのなかのひとつのインスタンスを識別する。
    * 実際のシステムでは、instance ラベルをテンプレート化し、マシンごとにすべてのネットワーク関連メトリクスを表示することになるだろう。
    ひとつのダッシュボードで複数のテンプレート変数を使うことさえできる。たとえば、Java のガベージコレクションダッシュボードは、一般に job のためにひとつ、instance のためにひとつ、そして使う Prometheus データソースの選択のためにひとつの変数を使うことになるだろう。
* [What is a job label for?](https://www.robustperception.io/what-is-a-job-label-for)
    * I would recommend is using the job label to organise applications which do the same thing, which almost always means processes running the same binary with the exact same configuration.
    * Something to avoid is encoding information beyond that what sort of binary/configuration the process is in a job label.

## exporter
* [one-agent-to-rule-them-all](https://www.robustperception.io/one-agent-to-rule-them-all)
    * 考え方として
* exporter down の検知
    * https://www.robustperception.io/alerting-on-down-instances
* instance は running のみで良いか
* group by をどうするか

## cAdvisor
* [docker.jsonnet](https://gist.github.com/rrreeeyyy/b276e5a319d9d0679e3b36af949e406c)
    * cAdvisor + ECS のサンプルとして
* [CFS Bandwidth Control](https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt)
    * 以下の CPU 周りのメトリクス理解に
* [](https://docs.signalfx.com/en/latest/integrations/agent/monitors/cadvisor.html)

### spec 系のメトリクス
* container_spec_cpu_period
* container_spec_cpu_shares
* container_spec_memory_limit_bytes
    * Memory limit for the container
* container_spec_memory_reservation_limit_bytes
* container_spec_memory_swap_limit_bytes

rate(container_cpu_usage_seconds_total[10m])/(container_spec_cpu_shares / container_spec_cpu_period)

sum by (instance, name) (rate(container_cpu_usage_seconds_total[60s]) * 60 * 1024) / on (instance, name) (container_spec_cpu_shares) / 60 * 100

sum by (instance, name) (rate(container_cpu_usage_seconds_total[60s]) * 60 * 1024) / on (instance, name) (container_spec_cpu_shares) / 60 * 100

sum by (instance, name, container_label_com_amazonaws_ecs_cluster) (rate(container_cpu_usage_seconds_total{name="ecs-web-ninnin-prd-4-web-ninnin-prd-86daf1b2fafaf1ba7700"}[60s]) * 60 * 1024) /  on (instance, name, container_label_com_amazonaws_ecs_cluster) (container_spec_cpu_shares{name="ecs-web-ninnin-prd-4-web-ninnin-prd-86daf1b2fafaf1ba7700"}) / 60 * 100
### ECS
* container_label_com_amazonaws_ecs_container_name
* container_label_com_amazonaws_ecs_task_definition_family
* container_label_com_amazonaws_ecs_cluster
* name

## metrics
### go client
* [EXPLORING PROMETHEUS GO CLIENT METRICS](https://povilasv.me/prometheus-go-metrics/)
* process_resident_memory_bytes
    * resident set memory size is number of memory pages the process has in real memory, with pagesize 4. This results in the amount of memory that belongs specifically to that process in bytes. This excludes swapped out memory pages.
* process_virtual_memory_bytes
    * virtual memory size is the amount of address space that a process is managing. This includes all types of memory, both in RAM and swapped out.
* process_open_fds
    * counts the number of files in /proc/PID/fd directory. This shows currently open regular files, sockets, pseudo terminals, etc.
* go_goroutines
    * calls out to runtime.NumGoroutine(), which computes the value based off the scheduler struct and global allglen variable.
* go_gc_duration_seconds
    * calls out to debug.ReadGCStats() with PauseQuantile set to 5, which returns us the minimum, 25%, 50%, 75%, and maximum pause times
* go_memstats_alloc_bytes
    * a metric which shows how much bytes of memory is allocated on the Heap for the Objects. The value is same as go_memstats_heap_alloc_bytes. This metric counts all reachable heap objects plus unreachable objects, GC has not yet freed.
* go_memstats_alloc_bytes_total
    * this metric increases as objects are allocated in the Heap, but doesn’t decrease when they are freed
    * Doing rate() on it should show us how many bytes/s of memory app consumes and is “durable” across restarts and scrape misses.
* go_memstats_heap_inuse_bytes
    * shows how many bytes in in-use spans
* go_memstats_stack_inuse_bytes
    * shows how many bytes of memory is used by stack memory spans, which have atleast one object in them. Go doc says, that stack memory spans can only be used for other stack spans, i.e. there is no mixing of heap objects and stack objects in one memory span.

## Alerts
* [What’s the difference between group_interval, group_wait, and repeat_interval?](https://www.robustperception.io/whats-the-difference-between-group_interval-group_wait-and-repeat_interval)
* [https://github.com/samber/awesome-prometheus-alerts](https://github.com/samber/awesome-prometheus-alerts)
* [When to Alert with Prometheus](https://www.robustperception.io/when-to-alert-with-prometheus)
* [Don’t put the value in alert labels](https://www.robustperception.io/dont-put-the-value-in-alert-labels)
    * つけるなら annotations につけないと、for がうまく機能しない
* [Absent Alerting for Jobs](https://www.robustperception.io/absent-alerting-for-jobs)
* [Alerting on approaching open file limits](https://www.robustperception.io/alerting-on-approaching-open-file-limits)

### alertmanager
* repeat_interval: 3h
    * How long to wait before re-sending a given alert that has already been sent.
    * まぁそんなにしつこくはしなくて良さそうだけど、severity で設定変えられるならそうすべきかと(出来るはず)

## Dashboard
* [Be discerning in what dashboards you share with users](https://www.robustperception.io/be-discerning-in-what-dashboards-you-share-with-users)
* [Avoid the Wall of Graphs](https://www.robustperception.io/avoid-the-wall-of-graphs)

## batch
* [Aggregating across batch job runs with push_time_seconds](https://www.robustperception.io/aggregating-across-batch-job-runs-with-push_time_seconds)

## 設定・運用
* [Limiting PromQL resource usage](https://www.robustperception.io/limiting-promql-resource-usage)
    * limit について参考になる
* [Unit testing rules with Prometheus](https://www.robustperception.io/unit-testing-rules-with-prometheus)
* [Unit testing alerts with Prometheus](https://www.robustperception.io/unit-testing-alerts-with-prometheus)
    * unit test
* [Monthly reporting with Prometheus and Python](https://www.robustperception.io/monthly-reporting-with-prometheus-and-python)
    * monthly みたいなレポートに PromQL は使いづらいので、その代案として
    * ちゃんと見るべき!
* [Extracting raw samples from Prometheus](https://www.robustperception.io/extracting-raw-samples-from-prometheus)
    * この場合は、query_range endpoint ではなく、query endpoint を使う！(上の記事もそうなっている)

## Reference
* [instrumenting-go-applications](https://banzaicloud.com/blog/instrumenting-go-applications/)
* [EXPLORING PROMETHEUS GO CLIENT METRICS](https://povilasv.me/prometheus-go-metrics/)
  * メトリクスの解説
* [How should pipelines be monitored?](https://www.robustperception.io/how-should-pipelines-be-monitored)