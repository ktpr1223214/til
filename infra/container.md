---
title: Docker
---

## 構築のプラクティス
* [コンテナ構築のおすすめの方法](https://cloud.google.com/solutions/best-practices-for-building-containers)

## network
### 前提知識
* [Docker/Kubernetesを扱う上で必要なネットワークの基礎知識（その1）](http://sagantaf.hatenablog.com/entry/2019/12/18/234553)
* [Docker/Kubernetesを扱う上で必要なネットワークの基礎知識（その2）](http://sagantaf.hatenablog.com/entry/2019/12/14/000948)

### 概要
* [Dockerのネットワークの仕組み](http://sagantaf.hatenablog.com/entry/2019/12/18/234553)
  * こちらを参考に
  * このページの外部と通信する時のネットワーク構成の図を外側の通信から書くと
    * （from internet）→ router で、private ip への変換（NAT）→ eth0 → iptables → docker0 → ... という理解

``` bash
# Docker のデフォルトでのネットワーク
$ docker network ls
# bridge/host/none が存在
# 何も指定しなければ bridge が使われる

# bridge の詳細を確認
$ docker network inspect bridge
```

#### 注意点
* docker0 birdge は macOS では存在しない（https://docs.docker.com/docker-for-mac/networking/#:~:text=There%20is%20no%20docker0%20bridge,actually%20within%20the%20virtual%20machine）
  * Because of the way networking is implemented in Docker Desktop for Mac, you cannot see a docker0 interface on the host. This interface is actually within the virtual machine.
  * macOS の IP は DHCP で取得することが多いはずなので、動的に変化するため

* macOS で
``` bash
# ゲートウェイとかを確認したい場合
$ netstat -rn
# もしくは、システム詳細設定のネットワークから、さらに詳細の TCP/IP とかからわかる
```

``` bash
# 適当なコンテナを動かしておく
$ docker run --rm -p 8080:80 nginx

# nginx container の ip を確認し、別コンテナから bridge 経由（docker0）でアクセスを試すことが可能
# macOS から見えないとしても、仕組みは同じはずなので
# docker ps -> docker inspect <container_id> で、割り当てられた ip を確認
# 別のコンテナ内部から curl で叩いてみるとちゃんと返ってくる
# ex. curl 172.17.0.2
$ curl <ip>

# en0 の ip 確認（private ip）
$ ifconfig
# host からでも、別の適当なコンテナ内部からでも、nginx の結果がちゃんと返る
$ curl <ip>:8080
# 後者の macOS host IP に対して、別のネットワークから叩きに行けるのはなぜ？
# ex. 試した場合だと、コンテナ ip が 172.17.0.2 一方で、macOS en0 の ip は、192.168.10.101 だった
# 何かしらを経由しないと後者にはたどり着けないはずだが、それはどうやって実現されているのか
# 適当なコンテナの中から traceroute をしてみると
$ traceroute 192.168.10.101
# 1  172.17.0.1 (172.17.0.1)  0.915 ms  0.845 ms  0.825 ms
# 2  192.168.10.101 (192.168.10.101)  3.682 ms  3.639 ms  3.613 ms
# つまり、コンテナの default gateway に向かって、そこから host に向かっている模様
```
* ちなみにコンテナで各種コマンドのインストールが必要かもしれない
``` bash
$ apt-get update && apt-get install traceroute iptables
```

### bridge
* [Use bridge networks](https://docs.docker.com/network/bridge/)
* User-defined の方が、default よりも優れている
  * default bridge network では、コンテナ間のやりとりは --link（considerd legacy）を使わない限りは IP アドレスでのみ
    * On a user-defined bridge network, containers can resolve each other by name or alias.
  * All containers without a --network specified, are attached to the default bridge network. This can be a risk, as unrelated stacks/services/containers are then able to communicate
  * などなど（上記ページを参照すること）

## storage

## Tips
* docker の alias で
