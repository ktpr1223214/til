---
title: ECS
---

## ECS

* ECS agent のインストール
``` bash
# cf. https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/ecs-agent-install.html#ecs-agent-install-al2
$ sudo amazon-linux-extras disable docker
$ sudo amazon-linux-extras install -y ecs; sudo systemctl enable --now ecs
```

* インストールしたインスタンスから AMI 作成→インスタンス作成などすれば、ECS cluster に入っていることが確認できるはず
    * userdata で

### user data
* [cloud-init-per ユーティリティ](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/bootstrap_container_instance.html#cloud-init-per)

## network mode
* [ECSでEC2インスタンスを利用する際のネットワークモードについて調べてみた](https://dev.classmethod.jp/etc/ecs-networking-mode/)

* bridge
    * ECSインスタンス（EC2） の任意のポートをコンテナのポートにマッピングして利用
    * ECSインスタンス（EC2） の ENI を複数のタスクが共有で利用
    * ENI を共用で利用するため SecurityGroup も共有   
    * ALB と組み合わせる場合は動的ポートで利用することが多い
    
* host
    * Docker の host
    * コンテナで export されたポートを ECSインスタンス(EC2)でも利用
        * そのため、一つのホストで同じポートは利用できない
        
* awsvpc
    * ENI がタスクごとにアタッチ
    * タスク間でのポートマッピングの考慮不要
    * ENI が独立しているため、ネットワークパフォーマンスの向上が見込める
    * ENI ごとに SecurityGroup を紐づけられる
    * ECS インスタンス本体とタスクで SecurityGroup を分けることも可能
    * VPC FlowLogs で観測可能
    * ALB と NLB に IP ターゲットとして登録が可能
    * ECSインスタンス（EC2）の ENI 上限には注意

## EFS との連携
* [EFS マウントヘルパー](https://docs.aws.amazon.com/ja_jp/efs/latest/ug/using-amazon-efs-utils.html#efs-mount-helper)
    * 基本はこれを使うべきっぽい

### 手動でマウントする場合
``` bash
# 事前確認(どちらでも)
$ mount | grep efs
$ df -Th

# /mnt/efs にマウントする場合
$ sudo mkdir /mnt/efs

# 起動時にマウントするように設定
# efs dns が fs-fea939df.efs.ap-northeast-1.amazonaws.com の場合
# TODO: amazon-efs-utils の場合を参照
$ echo 'fs-fea939df.efs.ap-northeast-1.amazonaws.com:/ /mnt/efs nfs4 nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 0 0' | sudo tee -a /etc/fstab

# -a /etc/fstabに記述されたファイルシステムを全てマウント
$ mount -a
```

* mount 手順
  * EFS の場合は、amazon-efs-utils を使ったこっちのほうが良さそう
``` bash
$ sudo yum install -y amazon-efs-utils
$ echo '<efs_dns>:/ <mount_path> efs defaults,_netdev 0 0' | sudo tee -a /etc/fstab
$ sudo mount -a -t efs defaults
```

ssh して
ps -ef で以下の時
root     19444 21433  0 03:07 ?        00:00:00 /usr/bin/docker-proxy -proto tcp -host-ip 0.0.0.0 -host-port 32777 -container-ip 172.17.0.2 -container-port 9090
root     19450 21483  0 03:07 ?        00:00:00 docker-containerd-shim -namespace moby -workdir /var/lib/docker/containerd/daemon/io.containerd.runtime.v1.linux/moby/34e3738cc12a3342d6b4fabecc581b3374975
nfsnobo+ 19488 19450  0 03:07 ?        00:00:00 /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_librarie

## service
* AWS console から service にいって、デプロイ・イベントなどから状況が確認できる

## service discovery
* IAM role について

* 
``` bash
# service_id は srv-... という形式で、console だと AWS Cloud Map から確認可能
$ aws servicediscovery list-instances --service-id <service_id> --region <region>

# namespace_name/service_name は AWS Cloud Map から。cluster_name は ECS cluster 名前
$ aws servicediscovery discover-instances --namespace-name <namespace_name> --service-name <service_name> --query-parameters ECS_CLUSTER_NAME=<cluster_name> --region ap-northeast-1

# 
# namespace_id は AWS Cloud Map から(ns-...)
$ aws servicediscovery get-namespace --id <namespace_id> --region <region>

# Route 53 ホストゾーン ID を使用して、ホストゾーンのリソースレコードセットを取得
$ aws route53 list-resource-record-sets --hosted-zone-id <hosted_zone_id> --region <region>

# SRV レコードを dig
# 帰ってきた値を更に dig ると A レコード(private ip アドレス)
$ dig SRV <domain_name>

# 更に SRV レコードを dig って返る port 番号 + private ip で health check とかサービスを叩ける
$ curl http://<private_ip>:<port>/<health_check_path>
```

### service discovery と ECS ネットワークモード
* [サービス検出 に関する考慮事項](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/service-discovery.html#service-discovery-considerations)
    * サービスタスクで指定されたタスク定義が awsvpc ネットワークモードを使用する場合、各サービスタスクで A レコードまたは SRV レコードを組み合わせて作成できます。SRV レコードを使用する場合、ポートが必要です。
    * サービスタスクで指定されたタスク定義が bridge または host ネットワークモードを使用する場合、SRV のレコードのみがサポートされる DNS レコードタイプです。各サービスタスクの SRV レコードを作成します。SRV レコードのコンテナ名とコンテナポートの組み合わせをタスク定義から指定する必要があります。
