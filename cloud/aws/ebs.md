---
title: EBS
---

## EBS
* 各インスタンスには、Amazon EBS ボリュームまたはインスタンスストアボリュームという、どちらかのルートデバイスボリュームが関連付けられている
* ブロックデバイスマッピングを使用することで、インスタンス起動時にアタッチする追加の EBS ボリューム or インスタンスストアボリュームを指定可
  * EBS ボリュームは実行中にアタッチすることもできるが、後者はブロックデバイスマッピングによる起動時アタッチのみ

## 関連コマンド
``` bash
# https://www.man7.org/linux/man-pages/man8/lsblk.8.html
# 使用可能なディスクデバイスとマウントポイント (該当する場合) を表示
# /dev/ プレフィックスは削除される
$ lsblk

# file -s: デバイスに関する情報 (ファイルシステムの種類など) を一覧表示
$ file -s /dev/xvdf

# ボリューム上にファイルシステムを作成
# 既にデータが入っているボリュームをマウントしていてこのコマンドを実行すると、既存データが削除されるので注意
# ex. /dev/xvdf をフォーマットして、/data にマウントするケース
$ sudo mkfs -t xfs /dev/xvdf
$ sudo mkdir /data
$ sudo mount /dev/xvdf /data
# アクセス許可の確認・設定

# 再起動後にも自動でマウントするには、/etc/fstab の設定が必要
# cf. https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/ebs-using-volumes.html#ebs-mount-after-reboot
```
