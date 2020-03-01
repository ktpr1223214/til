---
title: Node Exporter
---

## 基本
* [Linux Performance Analysis in 60,000 Milliseconds](https://netflixtechblog.com/linux-performance-analysis-in-60-000-milliseconds-accc10403c55)

### CPU
* user: The time spent in userland
* system: The time spent in the kernel
* iowait: Time spent waiting for I/O
* idle: Time the CPU had nothing to do
* irq&softirq: Time servicing interrupts
* guest: If you are running VMs, the CPU they use
* steal: If you are a VM, time other VMs "stole" from your CPUs

These modes are mutually exclusive. A high iowait means that you are disk or network bound, high user or system means that you are CPU bound.

* [CPU使用率は間違っている](https://yakst.com/ja/posts/4575)

### Memory
* [いまさら聞けないLinuxとメモリの基礎＆vmstatの詳しい使い方](https://qiita.com/kunihirotanaka/items/70d43d48757aea79de2d)
    * ページキャッシュ

### File System
* /dev/shm
    * tmpfs 専用のマウントポイント
    * tmpfs は仮想メモリベースのファイルシステム
* /run
    * 実行時の可変データ
* /sys/fs/cgroup
    * cgroupfs をマウントする先としての慣習
    * [LXCで学ぶコンテナ入門 －軽量仮想化環境を実現する技術](https://gihyo.jp/admin/serial/01/linux_containers/0003)
* /run/user/${uid}
    * There is a single base directory relative to which user-specific runtime files and other file objects should be placed.
        * [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
        * $XDG_RUNTIME_DIR で定義
* デバイスノードはシステムが起動するたびに生成されることになるので、 devtmpfs ファイルシステム上に保存
    * devtmpfs は仮想ファイルシステムで、メモリ上に置かれる
    * [7.3. デバイスとモジュールの扱いについて](http://lfsbookja.osdn.jp/7.6.ja/chapter07/udev.html)

### Disk
* node exporter は iostat -x の

### Network


### node での metrics
* node_memory_MemTotal_bytes
    * マシンの物理メモリの総容量
* node_memory_Cached_bytes
    * ページキャッシュ
* node_memory_Buffers_bytes
    * 書き込みバッファ
* node_memory_MemAvailable
    * カーネルが推計した本当に使えるメモリの容量
* node_memory_SwapTotal_bytes
    * total amount of swap space available
* node_filesystem_avail_bytes
    * ユーザが使えるディスクスペース
* node_filesystem_files
    * inode の数
* node_filesystem_files_free
    * フリー状態の inode の数
    * df -i と同じ情報
* node_network_receive_bytes_total/node_network_transmit_bytes_total
    * 送受信のネットワーク帯域幅利用
    * 8掛けで、bits per second に
* node_systemd_socket_accepted_connections_total
    * Total number of accepted socket connections
* node_systemd_units
    * Summary of systemd unit states

## サンプルの Grafana ダッシュボードを理解
* [Node Exporter Full](https://grafana.com/grafana/dashboards/1860)
    * 大人気っぽいコレ
    * Revision16
* ダッシュボードの上から解説

### Quick CPU/Mem/Disk
* CPU Busy
    * USE の U
    * (((count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu)))
    - avg(sum by (mode)(irate(node_cpu_seconds_total{mode='idle',instance=~"$node:$port",job=~"$job"}[5m])))) * 100)
    / count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu))
        * (count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu)))
            * コア数(cf. 下の Sys Load)
        * sum by (mode)(irate(node_cpu_seconds_total{mode='idle',instance=~"node-exporter:9100",job=~"node"}[5m]))
            * idle の node_cpu_seconds_total の秒あたり CPU 利用率を mode のみ残して総和(core とか)
            * これへの avg は下と同じ？
        * なので、分子は idle 以外に使った CPU 秒数で、それをコア数で割る * 100 なので、idle 以外での CPU 使用率(%)に

* Sys Load(5(15)m avg)
    * USE の S
    * avg(node_load5{instance=~"$node:$port",job=~"$job"}) / count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu)) * 100
        * avg(node_load5{instance=~"$node:$port",job=~"$job"})
            * 値は avg つけても変わらないと思うが、つけてないと割り算が出来無さそう？
        * count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu))
            * count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu)
                * これで、CPU 毎の count を取っている(ex. {cpu="0"} 8 {cpu="1"} 8)
            * ここは CPU コア数に等しい
            * cf. 入門 Prometheus 14.2.2.1
    * 割り算に100かけて%に
        * ex. 2CPU コア で、node_load5 の値が2なら、2/2*100 = 100%となる
        * LA はコア数基準で考えれば良いので

* RAM Used
    * USE の U
    * 100 -
    ((node_memory_MemAvailable_bytes{instance=~"$node:$port",job=~"$job"} * 100) / node_memory_MemTotal_bytes{instance=~"$node:$port",job=~"$job"})

* SWAP Used
    * ((node_memory_SwapTotal_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_SwapFree_bytes{instance=~"$node:$port",job=~"$job"}) 
    / (node_memory_SwapTotal_bytes{instance=~"$node:$port",job=~"$job"} )) * 100

* Root FS Used
    * 100 - ((node_filesystem_avail_bytes{instance=~"$node:$port",job=~"$job",mountpoint="/",fstype!="rootfs"} * 100) / node_filesystem_size_bytes{instance=~"$node:$port",job=~"$job",mountpoint="/",fstype!="rootfs"})
    
* CPU Cores
    * count(count(node_cpu_seconds_total{instance=~"$node:$port",job=~"$job"}) by (cpu))
    * 上述通り 

* RAM Total
    * node_memory_MemTotal_bytes{instance=~"$node:$port",job=~"$job"}
    * メトリクスまま

* SWAP Total
    * node_memory_SwapTotal_bytes{instance=~"$node:$port",job=~"$job"}

* Uptime
    * node_time_seconds{instance=~"$node:$port",job=~"$job"} - node_boot_time_seconds{instance=~"$node:$port",job=~"$job"}
    * 現在時間 - 起動時間?

### Basic CPU/Mem/Net/Disk
* CPU Basic
    * sum by (instance)(irate(node_cpu_seconds_total{mode="system",instance=~"$node:$port",job=~"$job"}[5m])) * 100
        * irate(node_cpu_seconds_total{mode="system",instance=~"node-exporter:9100",job=~"node"}[5m]) 
            * mode="system" の秒あたり CPU 利用率
            * sum by~で、instance 残して総和(CPU core)
    * 他も同様     
* Memory Basic
    * RAM Used
        * node_memory_MemTotal_bytes{instance=~"$node:$port",job=~"$job"}﻿- node_memory_MemFree_bytes{instance=~"$node:$port",job=~"$job"} -﻿(node_memory_Cached_bytes{instance=~"$node:$port",job=~"$job"}﻿+﻿node_memory_Buffers_bytes{instance=~"$node:$port",job﻿=~"$job"})
        * 一応利用可能なメモリとして、node_memory_MemFree_bytes だけでなく、ページキャッシュも書き込みバッファも回収して使えるから？
            * cf. 入門 Prometheus 7.5
* Network Traffic Basic
    * irate(node_network_receive_bytes_total{instance=~"$node:$port",job=~"$job"}[5m])*8
        * 上述通り、bit 単位にするために x8         
* Disk Space Used Basic
    * 下の Disk Space Used 参照

### CPU/Memory/Net/Disk
* CPU
    * sum by (mode)(irate(node_cpu_seconds_total{mode="system",instance=~"$node:$port",job=~"$job"}[5m])) * 100
    * まぁそのままで、コア数で割ったりしてないので100%を超えうる
* Memory Stack
    * 基本的にメトリクスまま
    * Apps
        * node_memory_MemTotal_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_MemFree_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_Buffers_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_Cached_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_Slab_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_PageTables_bytes{instance=~"$node:$port",job=~"$job"} - node_memory_SwapCached_bytes{instance=~"$node:$port",job=~"$job"}
        * これは諸々引いた感じかと    
* Network Traffic
    * これもまま 
* Disk Space Used
    * node_filesystem_size_bytes{instance=~"$node:$port",job=~"$job",device!~'rootfs'}
     - node_filesystem_avail_bytes{instance=~"$node:$port",job=~"$job",device!~'rootfs'}
        * rootfs は [ramfs, rootfs and initramfs](https://www.kernel.org/doc/Documentation/filesystems/ramfs-rootfs-initramfs.txt)       
* Disk IOps
    * irate(node_disk_reads_completed_total{instance=~"$node:$port",job=~"$job",device=~"[a-z]*[a-z]"}[5m])
    * device は [a-z] で始まって [a-z] で終わるもののみ
* I/O Usage Read/Write
    * irate(node_disk_read_bytes_total{instance=~"$node:$port",job=~"$job",device=~"[a-z]*[a-z]"}[5m])
* I/O Usage Times
    * irate(node_disk_io_time_seconds_total{instance=~"$node:$port",job=~"$job",device=~"[a-z]*[a-z]"} [5m])    

### Storage Disk
* 

### Storage Filesystem
* Filesystem space available
* File Nodes Free    
* File Descriptor        
* File Nodes Size
    * メトリクスまま
* Filesystem in ReadOnly
    * node_filesystem_readonly{instance=~"$node:$port",job=~"$job",device!~'rootfs'}
    * 0/1 っぽい         

### Network Traffic
* 

### Systemd
* Systemd Sockets
    * irate(node_systemd_socket_accepted_connections_total{instance=~"$node:$port",job=~"$job"}[5m])
    * [systemd/systemd-journald.socket](https://manpages.debian.org/jessie/systemd/systemd-journald.socket.8.en.html)
        * systemd-journald will forward all received log messages to the AF_UNIXSOCK_DGRAM socket /run/systemd/journal/syslog, if it exists, which may be used by Unix syslog daemons to process the data further.
    * 例えばだが、systemd での用途に使われるソケットのメトリクスかと
    * --collector.systemd が必要かと    
* Systemd Units State
    * node_systemd_units{instance=~"$node:$port",job=~"$job",state="activating"}
    * これは直感のままか
    * [systemd](https://www.freedesktop.org/software/systemd/man/systemd.html)
        * ここに state について書いてある
     
### Node Exporter
* node_scrape_collector_duration_seconds とかで、メトリクスまま

## Reference
* [kernel.org/doc/Documentation/filesystems/proc.txt](kernel.org/doc/Documentation/filesystems/proc.txt)
    * memory 系のメトリクスの情報とか
    * 必ず参照！
* [Understanding Machine CPU usage](https://www.robustperception.io/understanding-machine-cpu-usage)
* [Network interface metrics from the node exporter](https://www.robustperception.io/network-interface-metrics-from-the-node-exporter)
* [Filesystem metrics from the node exporter](https://www.robustperception.io/filesystem-metrics-from-the-node-exporter)
* [Mapping iostat to the node exporter’s node_disk_* metrics](https://www.robustperception.io/mapping-iostat-to-the-node-exporters-node_disk_-metrics)
* [Monitoring directory sizes with the Textfile Collector](https://www.robustperception.io/monitoring-directory-sizes-with-the-textfile-collector)