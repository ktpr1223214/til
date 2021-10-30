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
* node exporter は iostat -x と対応するメトリクスを持つ
  * [iostat](https://www.kernel.org/doc/Documentation/iostats.txt)
  * r/s: 1 秒の間にディスクドライバに発行された読み出し要求の数
    * 複数の IO をまとめようとするので、その話
    * rate(node_disk_reads_completed_total[5m])
  * w/s: 1 秒の間にディスクドライバに発行された書き込み要求の数
    * rate(node_disk_writes_completed_total[5m])
  * rrqm/s: 1 秒の間にドライバの要求キューにセットされ、マージされた読み出し要求の数
    * rate(node_disk_reads_merged_total[5m])
  * wrqm/s: 1 秒の間にドライバの要求キューにセットされ、マージされた書き込み要求の数
    * rate(node_disk_writes_merged_total[5m])
  * rkB/s: 1 秒の間にディスクドライバから読み出された kB 数
    * rate(node_disk_read_bytes_total[5m])（ただしこれは単位がバイトなので注意）
  * wkB/s: 1 秒の間にディスクドライバに書き込まれた kB 数
    * rate(node_disk_written_bytes_total[5m])（ただしこれは単位がバイトなので注意）
  * avgrq-sz: セクター内の平均要求サイズ（512 バイト）
    * avgrq-sz * 512 bytes/sector を計算すれば良い
    * マージされた後のサイズ
    * (rate(node_disk_read_bytes_total[5m]) + rate(node_disk_written_bytes_total[5m])) / (rate(node_disk_reads_completed_total[5m]) + rate(node_disk_writes_completed_total[5m]))
  * avgqu-sz: ドライバの要求で待機している要求とデバイスでアクティブに処理されている要求の合計の平均
    * rate(node_disk_io_time_weighted_seconds_total[5m])
  * await: I/O 応答時間の平均。ドライバの要求キューで待機している時間とデバイスの I/O の応答時間を含む（m 秒）
  * r_await: await と同じだが、読み出しのみ（m 秒）
    * rate(node_disk_read_time_seconds_total[5m]) / rate(node_disk_reads_completed_total[5m])
    * w_await も同様
    * await は、足し算してから割り算
      * (rate(node_disk_read_time_seconds_total[5m]) + rate(node_disk_write_time_seconds_total[5m])) / (rate(node_disk_reads_completed_total[5m]) + rate(node_disk_writes_completed_total[5m]))
  * svctm: ディスクデバイスの平均 I/O 応答時間（推定、m 秒）
  * %util: デバイスが I/O 要求を処理していてビジーだった時間の割合（使用率）
    * rate(node_disk_io_time_seconds_total[5m])
``` bash
# -x: 拡張統計（基本これ）-k: kB 数を使う -d: ディスクの報告 -z: アクティビティ 0 の集計を非表示
$ iostat -xkdz 1

Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await  svctm  %util
xvda              0.00   130.00   30.00   92.00  1384.00  6140.00   123.34     0.92  464.75   1.08  13.20
```

* node_disk_io_now
  * the number of IOs in progress
  * iostat では出力されない

### Network
* [Monitoring and Tuning the Linux Networking Stack: Receiving Data](https://blog.packagecloud.io/eng/2016/06/22/monitoring-tuning-linux-networking-stack-receiving-data/)

``` bash
$ netstat -s

Ip:
    176320801 total packets received
    178 with invalid addresses
    0 forwarded
    0 incoming packets discarded
    176320623 incoming packets delivered
    108229750 requests sent out
...
Tcp:
    169152 active connections openings
    306606 passive connection openings
    7710 SYNs to LISTEN sockets dropped
    7705 times the listen queue of a socket overflowed
    119 failed connection attempts
    1883 connection resets received
    429 connections established
    108475792 segments received
    113905327 segments send out
    1980 bad segments received.
    3679 segments retransmited
    TCPSynRetrans: 15841
...
Udp:
    67519821 packets received
    34503 packets to unknown port received.
    745432 packet receive errors
    154834 packets sent
    RcvbufErrors: 745432
...
TcpExt:

...
IpExt:
    InOctets: 127936695271
    OutOctets: 64332630877

```
* [netstatの統計情報を活用する](https://www.atmarkit.co.jp/fwin2k/win2ktips/311netstats/netstats.html)
* Ip
  * total packets received: 受信パケット総数
  * forwarded: 転送パケット数
    * 受信パケット総数に対する比率が高い場合、サーバーがパケットを転送するものか確認

* times the listen queue of a socket overflowed
  * node_netstat_TcpExt_ListenOverflows
* SYNs to LISTEN sockets dropped: SYN パケットの取りこぼし
  * node_netstat_TcpExt_ListenDrops
* TCPSynRetrans: number of SYN and SYN/ACK retransmits to break down retransmissions into SYN, fast-retransmits, timeout retransmits, etc.
  * node_netstat_TcpExt_TCPSynRetrans
* segments retransmited: 相手から受信確認が戻ってこないので、再送信したセグメントの総数
  * node_netstat_Tcp_RetransSegs
* bad segments received:
  * node_netstat_Tcp_InErrs

* Udp
  * packets received
    * node_netstat_Udp_InDatagrams
  * packet receive errors: UDP Datagrams that could not be delivered to an application
    * Incremented in several cases: no memory in the receive queue, when a bad checksum is seen, and if sk_add_backlog fails to add the datagram　らしい
    * node_netstat_Udp_InErrors
  * packets to unknown port received: UDP Datagrams received on a port with no listener
    * node_netstat_Udp_NoPorts
  * RcvbufErrors: Incremented when sock_queue_rcv_skb reports that no memory is available
    * ここに該当する node exporter のメトリクスは無さそう？
* IpExt
  * InOctets: 受信したオクテット数 = バイト数

* node_network_receive_packets_total
  * instance について、sum をとったものが、netstat -s の Ip: 配下の情報と対応している模様
* node_network_receive_errs_total
* node_network_transmit_queue_length

``` bash
$ netstat -i
Kernel Interface table
Iface       MTU Met    RX-OK RX-ERR RX-DRP RX-OVR    TX-OK TX-ERR TX-DRP TX-OVR Flg
eth0       9001   0 96374582      0      0      0 29500258      0      0      0 BMRU
lo        65536   0 77137912      0      0      0 77137912      0      0      0 LRU
```
* MTU: maximum transmission unit（byte 単位）
  * node_network_mtu_bytes
* OK: パケットの転送に成功
* ERR: パケットエラー
  * node_network_receive_errs_total
* DRP: パケットドロップ
  * node_network_receive_drop_total
* OVR: パケットオーバーラン

* なので、例えばエラー率を計算したい場合は
  * ```rate(node_network_transmit_errs_total[5m]) / rate(node_network_transmit_packets_total[5m])``` が使える
  * 注意: 分母は、送信失敗を含まない。が、実用上は気にすることはないか

``` bash
$ ip -s link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1
    ...
    RX: bytes  packets  errors  dropped overrun mcast
    18916327   164821   0       0       0       0
    TX: bytes  packets  errors  dropped carrier collsns
    18916327   164821   0       0       0       0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc pfifo_fast state UP mode DEFAULT group default qlen 1000
    ...
    RX: bytes  packets  errors  dropped overrun mcast
    20415321403 247193372 0       0       0       0
    TX: bytes  packets  errors  dropped carrier collsns
    42168916586 214789046 0       0       0       0
```

``` bash
$ sar -n EDEV

00時00分02秒     IFACE   rxerr/s   txerr/s    coll/s  rxdrop/s  txdrop/s  txcarr/s  rxfram/s  rxfifo/s  txfifo/s
00時10分01秒      eth0      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
00時10分01秒        lo      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00      0.00
```
* rxerr/s: 受信パケットエラー
* coll/s: コリジョン
  * node_network_transmit_colls_total
* rxdrop/s: ドロップ（バッファフル）した受信パケット
* rxfifo/s: FIFO オーバーランエラーを起こした受信パケット
  * node_network_receive_fifo_total（多分）

#### node での metrics
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

## USE Method and Node exporter
USE Method 観点で整理する Node exporter
### CPU
* Utilization: The time the CPU was busy (not in the idle thread)
* Saturation: The degree to which runnable threads are queued waiting their turn on-CPU
* Errors: CPU errors, including correctable errors

#### utilization
* rate(node_cpu_seconds_total{job="node-exporter"}[5m]) * 100 / ignoring (mode) count(count(node_cpu_seconds_total{job="node-exporter"}) by (cpu))
* CPU / Memory / Net / Disk -> CPU
* vmstat 1

Procs
  r: The number of processes waiting for run time.
  b: The number of processes in uninterruptible sleep.
Memory
  swpd: the amount of virtual memory used.
  free: the amount of idle memory.
  buff: the amount of memory used as buffers.
  cache: the amount of memory used as cache.
  inact: the amount of inactive memory. (-a option)
  active: the amount of active memory. (-a option)
Swap
  si: Amount of memory swapped in from disk (/s).
  so: Amount of memory swapped to disk (/s).
IO
  bi: Blocks received from a block device (blocks/s).
  bo: Blocks sent to a block device (blocks/s).
System
  in: The number of interrupts per second, including the clock.
  cs: The number of context switches per second.
CPU
  These are percentages of total CPU time.
  us: Time spent running non-kernel code. (user time, including nice time)
  sy: Time spent running kernel code. (system time)
  id: Time spent idle. Prior to Linux 2.5.41, this includes IO-wait time.
  wa: Time spent waiting for IO. Prior to Linux 2.5.41, included in idle.
  st: Time stolen from a virtual machine. Prior to Linux 2.6.11, unknown.

#### saturation
* la
  * Loosely, the load average is the number of processes running plus the those waiting to run
  * https://www.brendangregg.com/blog/2017-08-08/linux-load-averages.html
* node_load1, node_load5 and node_load15
* Quick CPU / Mem / Disk -> Sys Load
* uptime

#### error
* むずいっぽい

### memory
#### utilization
* node_memory_Cached_bytes
  * ページキャッシュ
* node_memory_Buffers_bytes
  * 書き込みバッファ

* ((node_memory_MemAvailable_bytes) / (node_memory_MemTotal_bytes) * 100)
* Basic CPU / Mem / Net / Disk ->
* free -b と比較した場合に
  * total: node_memory_MemTotal_bytes
  * used: node_memory_MemTotal_bytes - node_memory_MemFree_bytes - node_memory_Buffers_bytes - node_memory_Cached_bytes - node_memory_SReclaimable_bytes
  * shared: node_memory_Shmem_bytes
  * free: node_memory_MemFree_bytes
  * buff/cache: node_memory_Buffers_bytes + node_memory_Cached_bytes + node_memory_SReclaimable_bytes
    * Buffers: Relatively temporary storage for raw disk blocks shouldn't get tremendously large (20MB or so)
    * Cached: in-memory cache for files read from the disk (the pagecache). Doesn't include SwapCached
      * それぞれちょっとずれるが、多分ほぼ一致として良さげ
    * SReclaimable: Part of Slab, that might be reclaimed, such as caches
      * だから、slab 自体じゃなくてこっち（SReclaimable）なのは納得感あるが..
  * available: node_memory_MemAvailable_bytes
cf. https://stackoverflow.com/questions/59738885/prometheus-reports-a-different-value-for-node-memory-active-bytes-and-free-b
cf. https://milestone-of-se.nesuke.com/sv-basic/linux-basic/free-command/

理論的には、ページキャッシュ(node_memory_Cached_bytes)や書き込みバッファ(node_memory_Buffers_bytes)は回収して使えるが、そうすると一部のアプリケーションのパフォーマンスに悪影響を及ぼす。
さらに、カーネルには、スラブ(slab)やページテーブルなど、 メモリを使う部品が多数ある。

* https://access.redhat.com/solutions/406773
  * free -k and /proc/meminfo

* https://www.kernel.org/doc/Documentation/filesystems/proc.txt
  * これが網羅的

* vmstat 1 の buff + cache は free の buff/cache と一致
  * cache: node_memory_Cached_bytes + node_memory_SReclaimable_bytes
  * buff: node_memory_Buffers_bytes

https://qiita.com/kunihirotanaka/items/70d43d48757aea79de2d
ページキャッシュはファイルシステムに対するキャッシュであり、ファイル単位でアクセスするときに使用されるキャッシュです。
例えば、ファイルへデータを書き込んだときは、ページキャッシュにデータが残されるため、次回の読み込み時には HDD にアクセスすることなくデータを利用できます。 もうひとつのバッファキャッシュはブロックデバイスを直接アクセスするときに使用されるキャッシュです

#### saturation
* むずいっぽ
* pagein/pageout
* swapin/swapout
  * cf. https://superuser.com/questions/785447/what-is-the-exact-difference-between-the-parameters-pgpgin-pswpin-and-pswpou

* minor/major
  * https://en.wikipedia.org/wiki/Page_fault

pgpgin, pgpgout - number of pages that are read from disk and written to memory, you usually don't need to care that much about these numbers
pswpin, pswpout - you may want to track these numbers per time (via some monitoring like prometheus), if there are spikes it means system is heavily swapping and you have a problem

swapin/swapout は major fault の一種で、memory 不足の場合に発生するやつ、と考えれば良いはず
-> Node exporter のグラフで、major fault の方が swapout より少なかったりするのは何故か？
-> 恐らくだが、前者は操作の回数で後者は対象となったページ「数」だからかと（そう考えると変ではない）

swap 使ってなくても pageout が発生するのは？

Because of this, the term page-out means that a page was moved out of memory—which may or may not have included a write to a storage device

https://en.wikipedia.org/wiki/Memory_overcommitment

#### error
* やっぱむずいっぽい

### disk
When we are talking about disk utilization and saturation we have two additional dimensions to consider; disk capacity and disk throughput. Disk capacity is the amount of data that a disk can store (i.e. how full/empty is my disk) and disk throughput is how much can I read/write to the disk per second.

cf. https://brian-candler.medium.com/interpreting-prometheus-metrics-for-linux-disk-i-o-utilization-4db53dfedcfc

node_disk_reads_completed_total (field 1)
    This is the total number of reads completed successfully.
node_disk_reads_merged_total (field 2)
node_disk_writes_merged_total (field 6)
node_disk_discards_merged_total (field 13)
    Reads and writes which are adjacent to each other may be merged
    for efficiency.  Thus two 4K reads may become one 8K read before
    it is ultimately handed to the disk, and so it will be counted
    (and queued) as only one I/O.  This field lets you know how
    often this was done.
node_disk_read_bytes_total (field 3)
    This is the total number of bytes read successfully.
node_disk_read_time_seconds_total (field 4)
    This is the total number of seconds spent by all reads (as
    measured from __make_request() to end_that_request_last()).
node_disk_writes_completed_total (field 5)
    This is the total number of writes completed successfully.
node_disk_written_bytes_total (field 7)
    This is the total number of bytes written successfully.
node_disk_write_time_seconds_total (field 8)
    This is the total number of seconds spent by all writes (as
    measured from __make_request() to end_that_request_last()).
node_disk_io_now (field 9)
    The only field that should go to zero. Incremented as requests
    are given to appropriate struct request_queue and decremented as
    they finish.
node_disk_io_time_seconds_total (field 10)
    Number of seconds spent doing I/Os.
    This field increases so long as field 9 is nonzero.
node_disk_io_time_weighted_seconds_total (field 11)
    Weighted # of seconds spent doing I/Os.
    This field is incremented at each I/O start, I/O completion, I/O
    merge, or read of these stats by the number of I/Os in progress
    (field 9) times the number of seconds spent doing I/O since the
    last update of this field.  This can provide an easy measure of
    both I/O completion time and the backlog that may be
    accumulating.
node_disk_discards_completed_total (field 12)
    This is the total number of discards completed successfully.
node_disk_discarded_sectors_total (field 14)
    This is the total number of sectors discarded successfully.
node_disk_discard_time_seconds_total (field 15)
    This is the total number of seconds spent by all discards (as
    measured from __make_request() to end_that_request_last()).
node_disk_flush_requests_total (field 16)
    The total number of flush requests completed successfully.
node_disk_flush_requests_time_seconds_total (field 17)
    The total number of seconds spent by all flush requests.

#### utilization
##### capacity
* 1 - node_filesystem_avail_bytes / node_filesystem_size_bytes
* df -h
* Basic CPU / Mem / Net / Disk -> Disk Space Used Basic
##### throughput
* iostat -xz %util
* rate(node_disk_io_time_seconds_total[1m])
* Storage Disk -> Time Spent Doing I/Os

#### saturation
##### capacity
* capacity については特に該当なし

##### throughput
* "avgqu-sz" > 1, or high "await"
* iostat -xz
  * rrqm/s
    * rate(node_disk_reads_merged_total[*])
  * wrqm/s
    * rate(node_disk_writes_merged_total[*])
  * r/s: the number of reads per second calculated from the previous measurement iostat made (or since boot for the first one).
    * rate(node_disk_reads_completed_total[1m])
  * w/s
    * rate(node_disk_writes_completed_total[1m])
  * rkB/s
    * rate(node_disk_read_bytes_total[*])
  * wkB/s
    * rate(node_disk_written_bytes_total[*])
  * avgrq-sz
  * avgqu-sz
    * rate(node_disk_io_time_weighted_seconds_total[5m])
  * await
    * (rate(node_disk_read_time_seconds_total[5m]) + rate(node_disk_read_time_seconds_total[5m])) / (rate(node_disk_reads_completed_total[5m]) + rate(node_disk_writes_completed_total[5m]))
  * r_await
    * rate(node_disk_read_time_seconds_total[5m]) / rate(node_disk_reads_completed_total[5m])
  * w_await
  * svctm
    * deprecated
  * %util
    * rate(node_disk_io_time_seconds_total[1m])

r/s: The number (after merges) of read requests completed per second for the device.
w/s: The number (after merges) of write requests completed per second for the device.
rkB/s:The number of kilobytes¹ read from the device per second.
wkB/s: The number of kilobytes¹ written to the device per second.
rrqm/s: The number of read requests merged per second that were queued to the device.
wrqm/s: The number of write requests merged per second that were queued to the device.
%rrqm: The percentage of read requests merged together before being sent to the device.
%wrqm: The percentage of write requests merged together before being sent to the device.
r_await: The average time (in milliseconds) for read requests issued to the device to be served. This includes the time spent by the requests in queue and the time spent servicing them.
w_await: The average time (in milliseconds) for write requests issued to the device to be served. This includes the time spent by the requests in queue and the time spent servicing them.
aqu-sz: The average queue length of the requests that were issued to the device.
rareq-sz: The average size (in kilobytes) of the read requests that were issued to the device.
wareq-sz: The average size (in kilobytes) of the write requests that were issued to the device.
svctm: The average service time for I/O requests that were issued to the device. Unreliable, and has been removed from current versions of iostat.
%util: Percentage of elapsed time during which I/O requests were issued to the device (bandwidth utilization for the device). Device saturation occurs when this value is close to 100% for devices serving requests serially. But for devices serving requests in parallel, such as RAID arrays and modern SSDs, this number does not reflect their performance limits.

rrqm/s  マージされた1秒当たりの読み出し要求 ※※※
wrqm/s  マージされた1秒当たりの書き込み要求
r/s     1秒当たりに完了できた読み出し要求（マージ後）
w/s     1秒当たりに完了できた書き込み要求（マージ後）
rkB/s   1秒当たりの読み出しセクタ数
wkB/s   1秒当たりの書き込みセクタ数
avgrq-sz        1回で要求（ReQuest）された平均セクタサイズ（is the average size of each request, combining both read and write. It's calculated by iostat by dividing the bytes by the operations）
avgqu-sz        I/Oキュー（QUeue）の長さの平均
await   作成された要求が完了するまでの平均時間
r_await 作成された読み出し要求が完了するまでの平均時間
w_await 作成された書き込み要求が完了するまでの平均時間
svctm   デバイスに発行されたI/O要求の平均サービス時間（廃止予定）
％util  デバイスの帯域幅使用率
※※※同じデバイスに対する複数の読み出し（書き込み）要求を1回にまとめて処理した場合を示す値。

This is distinct from node_disk_io_now, which gives the instantaneous queue depth. Let’s say you are polling node_exporter every 5 seconds; then node_disk_io_now gives you the number of items which were in the queue at the sampling instant only. This value can vary massively from millisecond to millisecond, so the sampled value can be very noisy. This is the reason for the node_disk_io_time_weighted_seconds_total metric; by multiplying the queue depth by the amount of time the queue was at that given depth, and summing it up, the rate of increase of this metric gives you the average queue depth over the period in question.

* Storage Disk -> Disk IOs Weighted
#### error
* やっぱ..

### network
https://www.robustperception.io/network-interface-metrics-from-the-node-exporter

### utilization
* sar -n DEV 1 "rxKB/s"/max "txKB/s"/max
  * rxpck/s
  * txpck/s
  * rxkB/s
  * txkB/s
  * rxcmp/s
  * txcmp/s
  * rxmcst/s
* CPU / Memory / Net / Disk -> Network Traffic

### saturation
Saturation, without access to the underlying network capacity, will be hard to determine. We can use dropped packets as a proxy for network saturation:

* cat /proc/net/dev or ifconfig, "overruns", "dropped"
  * 前者の drop と後者の dropped は一致している
* rate(node_network_receive_drop_total[1m])
* rate(node_network_transmit_drop_total[1m])
* Network Traffic -> Network Traffic Drop

### err
* netstat -i or * cat /proc/net/dev or ifconfig, "overruns", "dropped"
  * RX-ERR/TX-ERR
* node_network_receive_errs_total and node_network_transmit_errs_total
* Network Traffic -> Network Traffic Errors

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
* Network Traffic By Packets
  * irate(node_network_receive_packets_total{instance="$node",job="$job"}[5m])
  * メトリクスまま
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
