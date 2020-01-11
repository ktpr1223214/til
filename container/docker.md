---
title: Docker
---

## Docker
### 仕組み
* Namespace
    * OS 起動時にはデフォルトの名前空間が存在し、プロセスはそこに属する
    * プロセスの起動時に、独立した名前空間を指定して実行すると、そのプロセスは別の名前空間で実行される

* Namespace の種類
    * マウント名前空間	 マウントの集合，操作
        * [マウント名前空間を適用する](https://www.ibm.com/developerworks/jp/linux/library/l-mount-namespaces.html)
    * UTS名前空間	 ホスト名，ドメイン名
    * PID名前空間	 プロセスID(PID)
    * IPC名前空間	 SysV IPCオブジェクト，POSIXメッセージキュー
    * ユーザ名前空間 UID，GID
    * ネットワーク名前空間 ネットワークデバイス，アドレス，ポート，ルーティングテーブル，フィルタなど

* 

* union filesystem
    * [第18回　Linuxカーネルのコンテナ機能 [7] ─ overlayfs](https://gihyo.jp/admin/serial/01/linux_containers/0018)

### Docker 関連パス
* /var/lib/docker/images/
    * Docker イメージ
* /var/lib/docker/containers/
    * Container 用イメージ層
* /var/lib/docker/volumes/
    * Volume

### 基本的なコマンド
``` bash
# docker ps -aq で container id のみ取得なので
$ docker rm $(docker ps -aq)
# 同様
# docker rmi $(docker images -aq)

# network/volume などの削除
$ docker container prune
$ docker network prune
$ docker image prune
$ docker volume prune

# Show docker disk usage
$ docker system df
# Remove unused data
$ docker system prune

# build の過程 下から上に
# <missing> とあるのは、別システムで build されており、local で not available という意味。無視で OK
$ docker history <image>

# container size on disk(https://docs.docker.com/storage/storagedriver/#container-size-on-disk)
# size: the amount of data (on disk) that is used for the writable layer of each container.
# virtual size: the amount of data used for the read-only image data used by the container plus the container’s writable layer size.
$ docker ps -s

# Docker for Mac の VM に connect
$ screen ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/tty
# ここから screen で適当にコマンド例
$ ls /var/lib/docker/
# ctrl a + d: detach ctrl a + k: kill
```


## Docker image
* layer とは
    * A Docker image is built up from a series of layers. 
    * Each layer represents an instruction in the image’s Dockerfile. 
    * Each layer except the very last one is read-only.

* storage driver とは
    * A storage driver handles the details about the way these layers interact with each other.

* 

### container と image
* 両者の違い
    * The major difference between a container and an image is the top writable layer.
    * All writes to the container that add new or modify existing data are stored in this writable layer.
    * When the container is deleted, the writable layer is also deleted. The underlying image remains unchanged.
* Docker における storage driver    
    * Docker uses storage drivers to manage the contents of the image layers and the writable container layer.
    * Each storage driver handles the implementation differently, but all drivers use stackable image layers and the copy-on-write (CoW) strategy.
* IO について    
    * Note: for write-heavy applications, you should not store the data in the container. 
    * Instead, use Docker volumes, which are independent of the running container and are designed to be efficient for I/O.

* [About storage drivers](https://docs.docker.com/storage/storagedriver/)

## Reference

* [Better docker image](https://speakerdeck.com/orisano/better-docker-image)
    * どのように速くするか・小さくするか