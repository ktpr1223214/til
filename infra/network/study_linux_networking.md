---
title
---

本やネットにある記事を参考にして、Linux にあるネットワーク周りの機能を学習する。
Linux サーバーの管理者として身に付けるべき知識には、Linux サーバーに特化しない IP ネットワークの一般的な知識と、Linux サーバーに固有の設定に関する知識
があるという。

## Linux ネットワーク周辺技術
* [LinuxのNetns/veth/Bridge/NATで仮想ネットワーク構築](https://blog.kamijin-fanta.info/2018/12/netns/)
  * ここの部分は、この記事冒頭から
* veth
  * Virtual Ethernet Device
  * 仮想的なネットワークインターフェイス
  * 2つのネットワークインターフェイスがペアになって機能（片方のネットワークインターフェイスにパケットが入ると、もう片方から出てくるという特性）
* iptables
  * ファイアウォールや NAT などの機能提供
  * 主に IP と TCP/UDP プロトコルに対する設定
* Bridge
  * L2 なブリッジインターフェースを作成する機能
  * 複数のネットワークインターフェースを、仮想的に L2 スイッチに繋げたような動作を行う
* Network Namespace
  * ルーティングテーブル・インターフェース等が分離された環境を複数作るための機構

## コマンド基本
* [ip COMMAND CHEAT SHEET for Red Hat Enterprise Linux](https://access.redhat.com/sites/default/files/attachments/rh_ip_command_cheatsheet_1214_jcs_print.pdf)
``` bash
# サーバーの ARP テーブルの内容を取得（ip neigh show の略）
$ ip n

# サーバー自身が持つ NIC の MAC アドレスの確認
$ ip l

# ネットワークインターフェイスに割り当てられている IP アドレスの確認
$ ip a

# ルーティングテーブルの確認
$ ip r
# 出力結果
# default via 192.168.1.1 dev eno1 proto static metric 100
# NIC eno1 から、ルーターの IP アドレス 192.168.1.1 にパケットを送信し、その後は先のネットワークに任せる。これがデフォルトゲートウェイ
# 192.168.1.0/24 dev eno1 proto kernel scope link src 192.168.1.11 metric 100
# 192.168.1.0/24 に向けたパケットは、NIC eno1 から送出される。scope link は同じサブネットなので、ルーターを介さないという意味
# via xx.xx.xx.xx: ネクストホップのルーター
# proto kernel: カーネルが自動生成した経路
# src: 送信元を指定

# EC2 の場合の例（subnet: 10.0.0.0/24）
# https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/VPC_Subnets.html
# CIDR ブロック 10.0.0.0/24 を持つサブネットの場合、10.0.0.1: VPC ルーター用に AWS で予約
# つまり、ネクストホップは 10.0.0.1 ルーターへ（これは、VPC ルーターといっているけど、つまりサブネットの単位でのルーターのはず）
# default via 10.0.0.1 dev eth0
# 10.0.0.0/24 dev eth0 proto kernel scope link src 10.0.0.7
# リンクローカルアドレス
# 169.254.169.254 dev eth0

# subnet: 10.0.1.0/24 の場合
# 10.0.2.0 が VPC ルーターで、つまりサブネット単位でのルーターかと
# default via 10.0.2.1 dev eth0
# 10.0.2.0/24 dev eth0 proto kernel scope link src 10.0.2.234
# 169.254.169.254 dev eth0
```
* EC2 の場合についてわかること
  * インターネットに出る（送信先がグローバルアドレス）場合は、10.0.0.1 でルーターに行き、
  そこで送信元のプライベート IPパブリックな IP に変換すると思えば良いか
  * 同じサブネットであれば、MAC アドレスで通信が完結すると思えば良く、そのルートが 10.0.0.0/24 dev eth0 proto kernel scope link src 10.0.0.7 かと
    * これが、VPC 内で通信を有効にするローカルルートということになる？
  * TODO: Security Group がどこで働くのかを理解したい
    * port とかも指定なので L4 のはずで（ICMP もあるので、L3 も一部か）、そこから何となくはわかるけどちゃんと確かめておきたい

## 実習
* 環境
  * Amazon Linux2 AMI EC2
    * ami-0a1c2ec61571737db

### NetworkManager
ネットワーク設定・管理をするもの
* setup
``` bash
$ sudo yum install -y NetworkManager
$ sudo systemctl start NetworkManager
```

* 基本コマンド
``` bash
# 定義済みのデバイス
$ nmcli d

# 定義済みの接続
$ nmcli c

# 特定のデバイスに関する設定
$ nmcli d show eth0

# 特定の接続に関する設定
$ nmcli c show eth0
```

### Namespace
* Namespace の作成
``` bash
# 作成 && 確認
# sudo ip netns add helloworld
$ sudo ip netns add <name>
$ ip netns list

# コマンド実行
# ex. sudo ip netns exec helloworld ip address show
$ sudo ip netns exec <name> <command>

$ ip netns add host1
$ ip netns add host2
$ ip netns ls

# veth 作成
# ex. sudo ip link add ns1-veth0 type veth peer name ns2-veth0
# 最初に書いた通り、veth はペアで
$ sudo ip link add <name_1> type veth peer name <name_2>

# veth を Network Namespace で使えるように
# ex. sudo ip link set ns1-veth0 netns ns1
# ex. sudo ip link set ns2-veth0 netns ns2
$ sudo ip link set <veth_name> netns <namespace>

# veth インターフェイスに IP アドレスを付与
$ sudo ip netns exec ns1 ip address add 192.0.2.1/24 dev ns1-veth0
$ sudo ip netns exec ns2 ip address add 192.0.2.2/24 dev ns2-veth0

# ネットワークインターフェイスを UP に
$ sudo ip netns exec ns1 ip link set ns1-veth0 up
$ sudo ip netns exec ns2 ip link set ns2-veth0 up

# 動作確認
$ sudo ip netns exec ns1 ping -c 3 192.0.2.2
# ここまでは同じネットワーク内部なので、ルーター不要
```

## Reference
* []()
* []()