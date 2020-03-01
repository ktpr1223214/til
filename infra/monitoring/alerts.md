---
title: Prometheus Alerts
---

## アラート例の解説
* [Awesome Prometheus alerts](https://awesome-prometheus-alerts.grep.to/rules) から必要そうなところを順に

### 1. Prometheus
* Prometheus not connected to alertmanager
    * prometheus_notifications_alertmanagers_discovered
    * そのままだけど、確かに大事っぽい
    * が、発火したときにどこに通知することに..?
* ExporterDown
    * up == 0
    * そのままで、とても大事かと

### 2. node-exporter
* Out of memory
    * node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10
    * メモリの空き割合が < 10%
* Unusual network throughput in(out)
    * sum by (instance) (irate(node_network_receive_bytes_total[2m])) / 1024 / 1024 > 100
        * 100MB という基準は何で？
        * 2m は何で？
* UnusualDiskReadRate(Write)
    * sum by (instance) (irate(node_disk_read_bytes_total[2m])) / 1024 / 1024 > 50 for 5m
        * .._disk_write_.. で Write
* OutOfDiskSpace
    * (node_filesystem_avail_bytes{mountpoint="/rootfs"}  * 100) / node_filesystem_size_bytes{mountpoint="/rootfs"} < 10 for 5m
        * mountpoint を / とか /var/lib/docker だとか色々変えて
* DiskWillFillIn4Hours
    * predict_linear(node_filesystem_free_bytes{fstype!~"tmpfs"}[1h], 4 * 3600) < 0 for 5m
    * 線形予測で、4時間後にファイル容量枯渇
* OutOfInodes
    * node_filesystem_files_free{mountpoint ="/rootfs"} / node_filesystem_files{mountpoint ="/rootfs"} * 100 < 10 for 5m
    * フリー状態の inode の割合
* UnusualDiskReadLatency(Write)
    * rate(node_disk_read_time_seconds_total[1m]) / rate(node_disk_reads_completed_total[1m]) > 100 for 5m
* HighCpuLoad
    * 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80 for 5m
        * irate より、rate のが良くはない？
* ContextSwitching
    * rate(node_context_switches_total[5m]) > 1000 for 5m
        * 1000ってのはアプリによる様子
        * 決められるのか？(https://github.com/samber/awesome-prometheus-alerts/issues/58)
* SwapIsFillingUp
    * (1 - (node_memory_SwapFree_bytes / node_memory_SwapTotal_bytes)) * 100 > 80 for 5m
* SystemdServiceCrashed
    * node_systemd_unit_state{state="failed"} == 1 for 5m
    * 起動時に --collector.systemd

### 3. cAdvisor
* Container killed
(sum(rate(container_cpu_usage_seconds_total[3m])) BY (container_label_com_amazonaws_ecs_container_name) * 100)

### 26. Blackbox
* ProbeFailed
    * probe_success == 0 for 5m
    * Blackbox exporter からの probe 失敗
* SlowProbe
    * avg_over_time(probe_duration_seconds[1m]) > 1 for 5m
    * 平均時間 > 1s
* HttpStatusCode
    * probe_http_status_code <= 199 OR probe_http_status_code >= 400 for 5m
    * HTTP Status Code が200-399じゃない
* SslCertificateWillExpireSoon
    * probe_ssl_earliest_cert_expiry - time() < 86400 * 30 for 5m
        * for 5m いるか？
    * SSL が 30日以内に expire
        * probe_ssl_earliest_cert_expiry - time() <= 0 なら既に
* HttpSlowRequests
    * avg_over_time(probe_http_duration_seconds[1m]) > 1 for 5m
* Slow ping
    * avg_over_time(probe_icmp_duration_seconds[1m]) > 1 for 5m