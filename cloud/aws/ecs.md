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

* stop-task
  * When StopTask is called on a task, the equivalent of docker stop is issued to the containers running in the task
  * This results in a SIGTERM value and a default 30-second timeout, after which the SIGKILL value is sent and the containers are forcibly stopped
  * If the container handles the SIGTERM value gracefully and exits within 30 seconds from receiving it, no SIGKILL value is sent
  * ECS container agent で ECS_CONTAINER_STOP_TIMEOUT を設定することで、上記 30 秒は変更可能
    * [Amazon ECS コンテナエージェントの設定](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/ecs-agent-config.html)

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
# port 番号と another_domain_name が返る
# 返ってきた another_domain_name を更に dig ると A レコード(private ip アドレス)
$ dig SRV <domain_name>
$ dig <another_domain_name>

# 更に SRV レコードを dig って返る port 番号 + private ip で health check とかサービスを叩ける
$ curl http://<private_ip>:<port>/<health_check_path>
```

* ECS service に紐づけて service discovery があるので、task のコンテナ1種類なら名前など指定せずとも検知してくれる
* service discovery で VPC 内の private DNS を見る場合には、SG で許可するのは別の SG or vpc_ip とかになり、public ip を間違って指定しないこと
* エフェメラルポートあたりも注意か

### service discovery と ECS ネットワークモード
* [サービス検出 に関する考慮事項](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/service-discovery.html#service-discovery-considerations)
    * サービスタスクで指定されたタスク定義が awsvpc ネットワークモードを使用する場合、各サービスタスクで A レコードまたは SRV レコードを組み合わせて作成できます。SRV レコードを使用する場合、ポートが必要です。
    * サービスタスクで指定されたタスク定義が bridge または host ネットワークモードを使用する場合、SRV のレコードのみがサポートされる DNS レコードタイプです。各サービスタスクの SRV レコードを作成します。SRV レコードのコンテナ名とコンテナポートの組み合わせをタスク定義から指定する必要があります。

## ECS 周辺の様々なロール
* AmazonEC2ContainerServiceforEC2Role
  * https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/instance_IAM_role.html
  * コンテナインスタンス内の ECS エージェントが ECS 呼び出しなどに必要な権限となる
    * なのでつけ先は EC2 instance
* AmazonEC2ContainerServiceRole
  * https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/service_IAM_role.html
  * Amazon ECS サービススケジューラは、ユーザーに代わって Amazon EC2 および Elastic Load Balancing の API を呼び出し、ロードバランサーを使ってコンテナインスタンスの登録および登録解除を行います。Amazon ECS サービスにロードバランサーをアタッチするには、それらのサービスの開始前に、使用する IAM ロールを作成する必要があります
    * なのでつけ先は ECS
* AmazonECSTaskExecutionRolePolicy
  * https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/task_execution_IAM_role.html
  * ECS task に紐付ける
    * task_role_arn - (Optional) The ARN of IAM role that allows your Amazon ECS container task to make calls to other AWS services.
      * 例えば S3 への権限とか、タスク側の権限
    * execution_role_arn - (Optional) The Amazon Resource Name (ARN) of the task execution role that the Amazon ECS container agent and the Docker daemon can assume.
      * この後者で紐付ける
      * こいつは ECS を実行する側の権限(ECR とか CWLogs とか)
  * EC2 cluster を使うのであれば、そちらに AmazonEC2ContainerServiceforEC2Role でも対応できるが、Task execution role で対応必要な場合もある
    * たとえば、ECS task の envirionment で、Parameter store を復号化して使う(secrets)の場合には、task execution role に該当権限が必要
      * [Amazon ECS シークレットで必須の IAM アクセス許可](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/task_execution_IAM_role.html#task-execution-secrets)

### 注意
* EC2 で ECS を動かす場合、task role と同じ権限は EC2 の IAM role に付与する必要はない
* AmazonEC2ContainerServiceforEC2Role があれば、ECR からの pull などを ECS task が行うことができる
* ただし、ECS task definition で task_execution_role を指定した場合には EC2 の IAM Role が、

### 参考
* [Confused by the role requirement of ECS](https://serverfault.com/questions/854413/confused-by-the-role-requirement-of-ecs)
  * 各種 role について説明があって良さそう
* [Amazon ECS の Amazon ECR エラー「CannotPullContainerError: API error」を解決する方法を教えてください](https://aws.amazon.com/jp/premiumsupport/knowledge-center/ecs-pull-container-api-error-ecr/)
  * ecs pull container error from ecr

## availability
* [Amazon ECS availability best practices](https://aws.amazon.com/jp/blogs/containers/amazon-ecs-availability-best-practices/)
    * これを参考に

## ECS + ELB
* target になるのは、あくまでインスタンス（port）
* デプロイでタスクが増える場合は、ELB からはターゲットが増えていることになる
  * ECS Task 増える -> ECS Service が TG に登録 -> draining して解除 -> Task 終了
* https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-target-groups.html#deregistration-delay

## サービスの更新
* [サービスの更新](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/update-service.html)

## デプロイ周辺
* デプロイメントは、タスク定義またはサービスの必要数を更新することでトリガーされる
  * デプロイ中、サービススケジューラは、最小ヘルス率と最大ヘルス率のパラメータを使用して、デプロイ戦略を判断
  * [その他のサービスの概念](https://docs.aws.amazon.com/ja_jp/AmazonECS/latest/developerguide/ecs_services.html#service_concepts)

* EC2+ASG+ECS な構成の場合に、インスタンス入れ替え
  * ASG からインスタンスを全部 detach -> このとき新しいインスタンスを起動するといった内容にチェック
  * このとき旧インスタンスは ASG の管理から外れるだけで ECS の管理から外れるわけではない
* 新しいインスタンスが起動したら、ECSで 1台ずつ Draining -> 新しいインスタンス上でタスクが起動
* 旧インスタンスを ECS のクラスターから deregistration する。これはコンテナインスタンスのIDのリンクから画面遷移して画面右上にボタンがある -> これで、ECS の管理から外れる
* 旧インスタンスを削除

### Draining
* クラスタからコンテナインスタンスを削除しなければならない場合 例:システム更新の実行、Dockerデーモンの更新、クラスタのスケールダウンサイズ がある
  * コンテナインスタンスのドレインにより、コンテナインスタンスを クラスター内のタスクに影響を与えずにクラスターを稼働させることが可能
* container instance を DRAINING にすると、新しい ECS Task が配置されなくなる
    * Service tasks on the draining container instance that are in the PENDING state are stopped immediately
    * If there are container instances in the cluster that are available, replacement service tasks are started on them
    * RUNNING の container instance は ECS Service の deployment parameter（minimumHealthyPercent・maximumPercent）に応じて停止・再配置が発生
      * If minimumHealthyPercent is below 100%, the scheduler can ignore desiredCount temporarily during task replacement
        * If tasks for services that do not use a load balancer are in the RUNNING state, they are considered healthy. Tasks for services that use a load balancer are considered healthy if they are in the RUNNING state and the container instance they are hosted on is reported as healthy by the load balancer
      * The maximumPercent parameter represents an upper limit on the number of running tasks during task replacement, which enables you to define the replacement batch size
        *  ex. desiredCount := 4 の場合、a maximum of 200% starts four new tasks before stopping the four tasks to be drained (provided that the cluster resources required to do this are available）
    * A container instance has completed draining when there are no more RUNNING tasks
* minimum healthy percent/maximum percent は update-service/container instance の DRAINING 時に参照される
  * [Updating a service](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/update-service.html)

### ECS update
* ECS Service + ELB の場合の update 処理流れ
  * When the service scheduler replaces a task during an update, the service first removes the task from the load balancer (if used) and waits for the connections to drain
  * Then, the equivalent of docker stop is issued to the containers running in the task. This results in a SIGTERM signal and a 30-second timeout, after which SIGKILL is sent and the containers are forcibly stopped
  * If the container handles the SIGTERM signal gracefully and exits within 30 seconds from receiving it, no SIGKILL signal is sent
  * The service scheduler starts and stops tasks as defined by your minimum healthy percent and maximum percent settings.

* [Does ECS update-service command marks the container instance state to draining when use with the --force-new-deployment option?](https://stackoverflow.com/questions/61409747/does-ecs-update-service-command-marks-the-container-instance-state-to-draining-w)
  1. ECS sends DeregisterTargets call and the targets change the status to "draining". New connections will not be served to these targets.
  2. If the deregistration delay time elapsed and there are still in-flight requests, the ALB will terminate them and clients will receive 5XX responses originated from the ALB.
  3. The targets are deregistered from the target group.
  4. ECS will send the stop call to the tasks and the ECS-agent will gracefully stop the containers (SIGTERM).
  5. If the containers are not stopped within the stop timeout period (ECS_CONTAINER_STOP_TIMEOUT by default 30s) they will force stopped (SIGKILL).

### ECS Stop task
* こいつを呼び出すと、ELB からの deregistration とかそういった ECS Service がやる作業は恐らく飛ばされるはずで断（502）が発生するはず
  * ECS Service のイベントをみると、service ~ has begun draining connections on 1 tasks と出たりはしてはいるが..
  * https://docs.aws.amazon.com/ja_jp/elasticloadbalancing/latest/application/load-balancer-troubleshooting.html#http-502-issues
* 例えば 2 つ ECS Task を動かしている場合に、1 つを落とすと上の事象。で、ECS Task を全部落とすと 最初は 502 でしばらくすると 503 に

### タスクの配置
* [Amazon ECS Task Placement](https://aws.amazon.com/jp/blogs/compute/amazon-ecs-task-placement/)
  * By default, ECS uses the following placement strategies:
    * When you run tasks with the RunTask API action, tasks are placed randomly in a cluster.
    * When you launch and terminate tasks with the CreateService API action, the service scheduler spreads the tasks across the Availability Zones (and the instances within the zones) in a cluster.

## ELB との構成
* ELB: public subnet
* ECS task: private subnet
  * ELB -> ECS は private ip によると思うので、ECS task が動くところ（EC2 とか）で、public ip は持たなくても処理を受けることは可能
  * ただし、コンテナからインターネットには出られないので、ECR なりを使うには VPC Endpoint が必須

## トラブルシューティング
* EC2 cluster なら、EC2 を探し出して /var/lib/docker/containers/ 辺りからログが見られるので、確認する(ls -l で時系列をみて)
  * ログ確認して想定されるログが出てない場合はアプリ側の出力方法を確認(ファイルのみに出力されているなど)
* ログは CloudWatch logs に出しておくなりすると見やすい(勿論 fluentd などで別に飛ばす場合はその限りでない)
* ALB を使った構成の場合は、ALB のログを S3 になりに出力しておくと、アプリへの接続がうまくいかない場合などのトラブルシューティングに役立つ
  * そもそも ALB まではいけているのか、など
* EC2 インスタンスが cluster に入らない場合
  * /etc/ecs/ecs.config に適切に cluster name が設定されているかどうか
  * tail /var/log/ecs/ecs-agent.log.... を確認し、エラーを吐いているか
  * EC2 instance は public ip を持つなど、ECS service endpoint との通信手段が必要なことに注意
    * cf. https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_instances.html#container_instance_concepts
* ECS service に関する IAM role は Fargate の場合付与しないので、間違った場合は消す・terraform init を再実行しないと、、でハマったこともあるので注意
* ECS task に関する良うわからんエラーは、定義の JSON を間違っている場合がある
