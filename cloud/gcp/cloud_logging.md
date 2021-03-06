---
title: Cloud Logging
---

## Cloud Logging
* Cloud 監査ログでは、Google Cloud のプロジェクト、フォルダ、組織ごとに管理アクティビティ、データアクセス、システム イベントの 3 つの監査ログが維持
* Google Cloud サービスによって、これらのログに監査ログエントリが書き込まれ、Google Cloud リソース内で「誰がどこでいつ何をしたか」を調べるのに役立つ
* [監査ログ付きの Google サービス](https://cloud.google.com/logging/docs/audit/services?hl=ja)

### ログの種類
* 管理アクティビティ監査ログ
  * 管理アクティビティ監査ログには、リソースの構成またはメタデータを変更する API 呼び出しやその他の管理アクションに関するログエントリが含まれる
  * これらのログは、たとえば、ユーザーが VM インスタンスを作成したときや Cloud Identity and Access Management 権限を変更したときに記録される
  * 監査ログは常に書き込まれ、構成したり無効にしたりすることはできない
  * 無料
* データアクセス監査ログ
  * データアクセス監査ログには、リソースの構成やメタデータを読み取る API 呼び出しや、ユーザー提供のリソースデータの作成、変更、読み取りを行うユーザー主導の API 呼び出しが含まれる
  * データアクセス監査ログは、非常に大きくなる可能性があるため、デフォルトで無効
  * 追加のコスト発生がありうる
* システム イベント監査ログ
  * リソースの構成を変更する Google Cloud 管理アクションのログエントリが含まれる
  * システム イベント監査ログは Google システムによって生成される
    * 直接的なユーザーのアクションによっては生成されない
  * システム イベント監査ログは常に書き込まれ、構成したり無効にしたりすることはできない
  * 無料

* 監査ログ名
```
projects/project-id/logs/cloudaudit.googleapis.com%2Factivity
projects/project-id/logs/cloudaudit.googleapis.com%2Fdata_access
projects/project-id/logs/cloudaudit.googleapis.com%2Fsystem_event

folders/folder-id/logs/cloudaudit.googleapis.com%2Factivity
folders/folder-id/logs/cloudaudit.googleapis.com%2Fdata_access
folders/folder-id/logs/cloudaudit.googleapis.com%2Fsystem_event

organizations/organization-id/logs/cloudaudit.googleapis.com%2Factivity
organizations/organization-id/logs/cloudaudit.googleapis.com%2Fdata_access
organizations/organization-id/logs/cloudaudit.googleapis.com%2Fsystem_event
```

``` bash
# ex
$ gcloud logging read "logName : projects/<project-name>/logs/cloudaudit.googleapis.com%2Factivity"
```

## エクスポート
* エクスポートは、エクスポートするログエントリを選択するクエリを書くことと、Cloud Storage、BigQuery、Pub/Sub でエクスポート先を選択することを含む
  * クエリとエクスポート先は、シンクと呼ばれるオブジェクトに保持される
* Google Cloud のプロジェクト、組織、フォルダ、請求先アカウントでシンクを作成可
* シンクの構成
  * シンク識別子: シンクの名前
  * 親リソース: シンクを作成するリソース
    * プロジェクトや組織・フォルダなど
    * シンクは、その親リソースに属するログのみをエクスポート可能
    * 例外の集約エクスポートについては後述
  * ログフィルタ: このシンクからエクスポートするログエントリ
  * エクスポート先: クエリに一致するログエントリを送信する単一の場所
    * ログはどのプロジェクト内のエクスポート先にもエクスポート可能だが、エクスポート先がシンクのサービス アカウントをライターとして承認する必要あり
  * ライター ID: サービス アカウント名。エクスポート先のオーナーは、このサービス アカウントにエクスポート先への書き込み権限を付与する必要がある
* シンクの仕組み
  * ログエントリがプロジェクト、フォルダ、請求先アカウント、または組織のリソースに到着するたびに、Logging はログエントリをそのリソース内のシンクと比較
  * クエリがログエントリと一致する各シンクは、ログエントリのコピーをシンクのエクスポート先に書き込み
  * 新しいログエントリに対してのみエクスポートが行われるため、シンクが作成される前に Logging が受信したログエントリはエクスポートできない
* 集約シンク
  * Google Cloud 組織のすべてのプロジェクト、フォルダ、請求先アカウントからログエントリをエクスポートできる集約シンクを作成できる
  * たとえば、組織のプロジェクトから 1 か所に集中して、監査ログエントリを集約、エクスポート可能
  * 集約シンク機能を使用するには、Google Cloud の組織またはフォルダにシンクを作成し、シンクの includeChildren パラメータを True に設定
    * これにより、そのシンクは組織またはフォルダに加えて、それに含まれているすべてのフォルダ、請求先アカウント、プロジェクトから（再帰的に）ログエントリをエクスポートできるようになり、シンクのクエリを使用して、プロジェクト、リソースタイプ、名前のついたログからログエントリを指定できる
* エクスポートしたログの可用性
  * Cloud Storage バケットへのログエントリの保存は 1 時間ごとに一括して行われます。最初のエントリが表示されるまでに 2～3 時間かかることがあります。
* 時間例
  * 18:13: service account 作成
  * 19:06: gcs に連携
    * 09:00:00_09:59:59_S0.json
    * JST で 18:00~18:59

* エクスポートの詳細
  * [Using exported logs](https://cloud.google.com/logging/docs/export/using_exported_logs)
* Cloud Storage の場合
  * Cloud Storage バケットへのログエントリの保存は 1 時間ごとに一括
  * ```my-gcs-bucket/\<LogEntry\>/YYYY/MM/DD/```
    * LogEntry は https://cloud.google.com/logging/docs/export/using_exported_logs#the_logentry_type
    * syslog/request_log/activity
  * ファイル内のログエントリの並べ替え順は一定ではなく、保証されないことに注意!!!

* ログの解釈  
  * [サンプル監査ログエントリの解釈](https://cloud.google.com/logging/docs/audit/understanding-audit-logs#interpreting_the_sample_audit_log_entry)  
  * [LogEntry](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry)
    * An individual entry in a log.    
  * [AuditLog](https://cloud.google.com/logging/docs/reference/audit/auditlog/rest/Shared.Types/AuditLog)
    * Common audit log format for Google Cloud Platform API operations.        
  * [AuditData](https://cloud.google.com/iam/docs/reference/rest/v1/AuditData)
    * Audit log information specific to Cloud IAM admin APIs. This message is serialized as an Any type in the ServiceData message of an AuditLog message.

* 巨大 or 処理時間がかかる場合の audit log
  * オペレーションが非同期で実行される、またはサイズの大きい AuditLog レコードが生成されると、1 つの監査対象オペレーションが複数のログエントリに分割
  * 同じオペレーションに複数のログエントリが存在する場合、LogEntry オブジェクトに operation フィールドが含まれ、同じオペレーションのエントリの LogEntry.operation.id と LogEntry.operation.producer に同じ値が設定
