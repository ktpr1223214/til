---
title: EC2
---

## EC2
* インスタンスを起動するときは、ルートデバイスボリュームに格納されているイメージを使用してインスタンスがブート
  * Amazon EC2 のサービス開始当初は、すべての AMI が「Amazon EC2 インスタンスストア backed」だった
    * つまり、AMI から起動されるインスタンスのルートデバイスは、Amazon S3 に格納されたテンプレートから作成されるインスタンスストアボリューム
  * Amazon EBS の導入後は Amazon EBS を基にした AMI も導入
    * つまり、AMI から起動されるインスタンスのルートデバイスが、Amazon EBS スナップショットから作成される Amazon EBS ボリューム
  * 後者の利用が推奨されている
    * 起動が高速で、永続化ボリュームを使用しているから
* Amazon EBS-Backed と Instance Store-Backed の違い
  * [ルートデバイスのストレージ](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/ComponentsAMIs.html#storage-for-the-root-device)を参照
* [Amazon EC2 ルートデバイスボリューム](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/RootDeviceStorage.html)
  * ルートデバイスの詳細はここで
* [インスタンスメタデータの取得](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html)
``` bash
# インスタンスに割り当てられているグローバル IP の取得
$ curl -s http://169.254.169.254/latest/meta-data/public-ipv4
```

## AMI
* [Linux AMI 仮想化タイプ](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/virtualization_types.html)
  * 準仮想化 (PV) およびハードウェア仮想マシン (HVM) のいずれか
    * 主な違いは、起動の方法と、パフォーマンス向上のための特別なハードウェア拡張機能 (CPU、ネットワーク、ストレージ) を利用できるかどうか
  * 最適なパフォーマンスを得るために、インスタンスを起動するときには、現行世代のインスタンスタイプと HVM AMI を使用することを推奨

## IP
* パブリック IP アドレスは、ネットワークアドレス変換 (NAT) によって、プライマリプライベート IP アドレスにマッピングされる

## Security Group
* [覚えておくべき4種のファイアウォール](https://www.internetacademy.jp/it/management/security/four-types-of-firewalls.html)
  * これでいうと、SG はダイナミックパケットフィルタリングと思えば良いか

## 起動・停止・再起動・終了
* stop
  * [インスタンスの停止と起動](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/Stop_Start.html#instance_stop)
  * 停止できるのは Amazon EBS-Backed インスタンスのみ
    * i-04cba66b82495dbee
    * 低指定、
  * ホストコンピュータの RAM またはホストコンピュータのインスタンスストアボリュームに保存されたデータはなくなる
  * 殆どの場合、インスタンスは基盤となる新しいホストコンピュータが起動したときに移行
  * インスタンスが Auto Scaling グループにある場合、Amazon EC2 Auto Scaling サービスはインスタンスを異常と判断して停止し、場合によってはそれを終了して代わりのインスタンスを起動
  * 停止中に変更できる属性
    * インスタンスタイプ
    * ユーザーデータ
    * Kernel
    * RAM ディスク

* terminate
  * インスタンスが終了すると、Amazon EC2 はアタッチされた各 Amazon EBS ボリュームの DeleteOnTermination 属性の値を使用して、ボリュームを保持するか削除するかを決定
    * ルートボリューム
      * デフォルトでは、インスタンスのルートボリュームの DeleteOnTermination 属性は true
    * ルート以外のボリューム
      * デフォルトでは、インスタンスにルート以外の EBS ボリュームをアタッチするときは、DeleteOnTermination 属性が false に設定

* reboot
  * インスタンスを再起動すると、インスタンスのパブリック DNS 名 (IPv4)、プライベート IPv4 アドレス、IPv6 アドレス (該当する場合)、およびインスタンスストアボリューム上のすべてのデータが保持
    * インスタンスは同じホストコンピュータに残る
  * インスタンスの再起動は、オペレーティングシステムの再起動と同等

* stop-start
  * インスタンスのステータスチェックに失敗するか、インスタンスでアプリケーションが想定通りに動作しておらず、インスタンスのルートボリュームが Amazon EBS である場合、インスタンスの停止と起動を行い、問題が解決するか試してみることが可能
  * インスタンスを起動すると、pending状態に移行し、ほとんどの場合は新しいホストコンピュータに移動される(ホストコンピュータに問題がない場合、インスタンスは同じホストコンピュータに残る可能性があります)
  * インスタンスの停止と起動を行うと、前のホストコンピュータ上のインスタンスストアボリューム上に存在していたすべてのデータが失われる

* 違い
  * [参考](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html#lifecycle-differences)
  * ホストコンピュータ
    * 再起動: インスタンスは、同じホストコンピュータで保持
    * 停止/開始 (Amazon EBS-Backed インスタンスのみ): ほとんどの場合、インスタンスは新しいホストコンピュータに移動(ホストコンピュータに問題がない場合、インスタンスは同じホストコンピュータに残る可能性あり)
    * 休止 (Amazon EBS Backed インスタンスのみ）: 同上
    * 終了: なし
  * プライベート IPv4 アドレスとパブリック IPv4 アドレス
    * 再起動: 同一のまま保持
    * 停止/開始 (Amazon EBS-Backed インスタンスのみ): インスタンスはプライベート IPv4 アドレスを保持。インスタンスは、Elastic IP アドレス (停止/起動の際に変更されない) を持っていない限り、新しいパブリック IPv4 アドレスを取得
    * 休止 (Amazon EBS Backed インスタンスのみ）: 同上
    * 終了: なし
  * Elastic IP アドレス (IPv4)
    * 終了時は関連付けが解除され、それ以外は関連づけられたまま維持
  * IPv6 アドレス
    * 再起動: アドレスは同一のまま保持
    * 停止/開始 (Amazon EBS-Backed インスタンスのみ): インスタンスは、IPv6 アドレスを保持する
    * 休止 (Amazon EBS Backed インスタンスのみ）: 同上
    * 終了: なし
  * インスタンスストアボリューム
    * 再起動: データは保持される
    * 停止/開始 (Amazon EBS-Backed インスタンスのみ): データは消去される
    * 休止 (Amazon EBS Backed インスタンスのみ）: 同上
    * 終了: データは消去される
  * ルートデバイスボリューム
    * 再起動: ボリュームは保持される
    * 停止/開始 (Amazon EBS-Backed インスタンスのみ): 同上
    * 休止 (Amazon EBS Backed インスタンスのみ）: 同上
    * 終了: ボリュームはデフォルトで削除される

## ステータスチェック
* システムステータスのチェック
  * システムステータスチェックが失敗した場合、AWS が問題を解決するのを待つか、自分で解決できるかを選択可能
  * インスタンスが実行されているAWSシステムを監視Amazon EBS でバックアップされたインスタンスの場合は、インスタンスを自分で停止および起動することができる。通常、インスタンスは新しいホストに移行される
  * 原因
    * ネットワーク接続の喪失
    * システム電源の喪失
    * 物理ホストのソフトウェアの問題
    * ネットワーク到達可能性に影響する、物理ホスト上のハードウェアの問題
* インスタンスステータスのチェック
  * 個々のインスタンスのソフトウェアとネットワークの設定をモニタリング
    * Amazon EC2 は、ネットワークインターフェイス (NIC) にアドレス解決プロトコル (ARP) リクエストを送信することによってインスタンスの健全性をチェック
  * これらのチェックでは、ユーザーが関与して修復する必要のある問題が検出される
    * インスタンスステータスチェックが失敗した場合は通常、自分自身で (たとえば、インスタンスを再起動する、インスタンス設定を変更するなどによって) 問題に対処する必要がある
  * 原因
    * 失敗したシステムステータスチェック
    * 正しくないネットワークまたは起動設定
    * メモリの枯渇
    * 破損したファイルシステム
    * 互換性のないカーネル

* 意図的にステータスチェックを失敗させたい場合インスタンスステータスは以下で可能
``` bash
$ ifconfig eth0 down
```

## ストレージ
* EC2 がサポートするデータストレージオプション
  * EBS
  * Amazon EC2 インスタンスストア
  * EFS
  * S3
* EBS
  * Amazon EBS は、実行中のインスタンスにアタッチできる、堅牢なブロックレベルのストレージボリューム
  * Amazon EBS は、細かな更新を頻繁に行う必要があるデータを対象とした主要ストレージデバイスとして使用可能
    * ex. インスタンスでデータベースを実行するときなど
  * EBS ボリュームは、1 つのインスタンスにアタッチできる、未加工、未フォーマットの外部ブロックデバイスのように動作
  * データのバックアップコピーを保持するには、EBS ボリュームのスナップショットを作成して Amazon S3 に保存
* Amazon EC2 インスタンスストア
  * 多くのインスタンスは、ホストコンピュータに物理的にアタッチされたディスクからストレージにアクセス可能
    * このディスクストレージは、インスタンスストアと呼ばれる
  * インスタンスストアは、インスタンス用のブロックレベルの一時ストレージを提供。インスタンスストアボリュームのデータは、関連するインスタンスの存続中にのみ保持され、インスタンスを停止または終了すると、インスタンスストアボリュームのすべてのデータが失われる
* EFS/S3
  * まぁはい

### Amazon EC2 インスタンスストア
* インスタンスストアは、インスタンス用のブロックレベルの一時ストレージを提供
* このストレージは、ホストコンピュータに物理的にアタッチされたディスク上にある
* インスタンスストアは、頻繁に変更される情報 (バッファ、キャッシュ、スクラッチデータ、その他の一時コンテンツなど) の一時ストレージに最適
* インスタンスストアは、ブロックデバイスとして表示される 1 つ以上のインスタンスストアボリュームで構成
  * インスタンスストアボリュームの仮想デバイスは ephemeral\[0-23\]

* インスタンスストアボリュームのインスタンスブロックデバイスマッピングの表示
``` bash
$ curl http://169.254.169.254/latest/meta-data/block-device-mapping/
```

### EBS
* Amazon Elastic Block Store (Amazon EBS) は、EC2 インスタンスで使用するためのブロックレベルのストレージボリュームを提供
  * EBS ボリュームの動作は、未初期化のブロックデバイスに似ている
* インスタンスにアタッチされているボリュームの設定は動的に変更可能
* インスタンスにボリュームをアタッチする場合の、ボリュームのデバイス名について
  * [Linux インスタンスでのデバイスの名前付け](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/device_naming.html)
* ボリュームを使えるようにする手順
  * [Linux で Amazon EBS ボリュームを使用できるようにする](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/ebs-using-volumes.html)
``` bash
# 現在利用できるブロックデバイスを一覧表示（/dev/ プレフィックスは除外される）
$ lsblk

# file: determine file type
# Normally, file only attempts to read and determine the type of argument files which stat(2) reports are ordinary files.  This prevents problems, because reading special files may have peculiar consequences.  Specifying the -s option causes file to also read argument files which are block or character special files.  This is useful for determining the filesystem types of the data in raw disk partitions, which are block special files.
# /dev/xvdb は 1 例
$ sudo file -s /dev/xvdb

# ファイルシステム実構築の場合のみ
# mkfs: build a Linux filesystem
# -t: Specify the type of filesystem to be built
$ sudo mkfs -t xfs /dev/xvdb

# ボリュームのマウントポイントディレクトリを作成
$ sudo mkdir /data

# もし、/data に何かを先に書き込んでいて場合は、マウント時に見えなくなる
# マウント
$ sudo mount /dev/xvdb /data
# 確認（どっちでも）
$ df -Th
$ lsblk

# アンマウント（書き込んでいたデータは戻る）
$ sudo umount /data

# マウントを永続化したい場合
# /dev/xvdb /data xfs defaults,nofail 0 0
# などを書き込む
$ sudo vi /etc/fstab
$ sudo unmout /data
$ sudo mount -a
# うまくいっていれば、マウントされる
```
* [ボリュームサイズ変更後の Linux ファイルシステムの拡張](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/recognize-expanded-volume-linux.html)

### ブロックデバイスマッピング
* ブロックデバイスは、一連のバイトまたはビット (ブロック) でデータを移動するストレージデバイス
* これらのデバイスはランダムアクセスをサポートし、通常は、バッファされた I/O を使用
  * たとえば、ハードディスク、CD-ROM ドライブ、フラッシュドライブなどがブロックデバイス
* EC2 がサポートするブロックデバイスは 2 種類
  * インスタンスストアボリューム (基盤となるハードウェアがインスタンスのホストコンピュータに物理的にアタッチされている仮想デバイス)
  * EBS ボリューム (リモートストレージデバイス)
* ブロックデバイスマッピングでは、インスタンスにアタッチするブロックデバイス (インスタンスストアボリュームと EBS ボリューム) を定義
  * ブロックデバイスマッピングは、AMI 作成プロセスの一環として、AMI から起動されるすべてのインスタンスによって使用されるように指定可能
  * また、インスタンスの起動時にブロックデバイスマッピングを指定することもできる
    * AMI のマッピングはこの場合上書きされる

## ASG

* Auto Scaling インスタンスのヘルスチェック
  * Auto Scaling インスタンスのヘルスステータスは、正常または異常のどちらか
  * Amazon EC2 Auto Scaling がインスタンスを異常ありとマークすると、置き換えがスケジュールされる
  * インスタンスに異常があるという通知は、Amazon EC2、Elastic Load Balancing (ELB)、カスタムヘルスチェックのうち 1 つ以上のソースから送られる可能性がある
    * インスタンスを損なう可能性があるハードウェアとソフトウェアの問題を特定するために Amazon EC2 によって提供されるステータスチェック。Auto Scaling グループのデフォルトのヘルスチェックは EC2 ステータスチェックのみ
    * Elastic Load Balancing (ELB) が提供するヘルスチェック。このヘルスチェックはデフォルトで無効になっていますが、有効にできる
    * カスタムヘルスチェック
