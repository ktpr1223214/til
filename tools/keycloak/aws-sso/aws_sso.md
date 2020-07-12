---
title: AWS SSO by Keycloak
---

## 概要
### 登場人物
* AWS
  * SP（Service Provider）
    * 有効なアサーションを提示された場合に、対象ユーザーの属性や属性値でコンテンツを提供
* Keycloak
  * IdP（ID Provider）
    * ユーザー認証成功後に、アサーションを発行
      * アサーション: 対象とする主体の認証や属性あるいは資源に関する認可権限の証明
* IAM SAML 2.0 ID プロバイダー
  * SAML 2.0 基準をサポートする外部 ID プロバイダー (IdP) を記述する IAM のエンティティ
  * 今回だと、Keycloak のことを IAM エンティティとして扱うためのものかと。IAM の SAML プロバイダーは IAM 信頼ポリシーでプリンシパルとして使用される

### 設定の流れ
* SAML プロバイダーを作成
  * SAML 2.0 基準をサポートする外部 ID プロバイダー (IdP) を記述する IAM のエンティティとして
* IAM ロールを作成する
  * ロールは AWS のアイデンティティであり、それ自体には (ユーザーのような) 認証情報がない
  * しかし、この例で、ロールが動的に割り当てられるフェデレーティッドユーザーは、組織の IdP から認証される
  * このロールで、組織の IdP が AWS にアクセスするための一時的なセキュリティ認証情報をリクエストできるようにする
  * ロールに割り当てられているポリシーは、フェデレーティッドユーザーが AWS で実行できることを決定
* AWS およびフェデレーティッドユーザーが使用するロールに関する情報で IdP を設定し、SAML の信頼を完了
  * これを、IdP と AWS 間の証明書利用者の設定という

### SAML プロバイダーの作成
* IAM ID プロバイダーを作成するには、IdP から SAML メタデータドキュメントを取得する必要がある
  * このドキュメントには、発行者名、失効情報、およびキーが含まれており、これらを使用して IdP から受け取った SAML 認証レスポンス (アサーション) を検証できる
  * [SAML V2.0 Metadata Guide](https://www.oasis-open.org/committees/download.php/51890/SAML%20MD%20simplified%20overview.pdf)

### 証明書利用者の信頼およびクレームの追加によって SAML 2.0 IdP を設定
* IdP に対し、サービスプロバイダーとしての AWS について通知
  * IdP と AWS の関係に対する証明書利用者の信頼の追加と呼ばれる
* IdP によって設定方法は異なるが ```https://signin.aws.amazon.com/static/saml-metadata.xml``` を利用する
* また、IdP で、AWS を証明書利用者として指定する適切なクレームルールを作成する必要がある
  * IdP が AWS エンドポイントに対して SAML レスポンスを送信する場合、これには 1 つ以上のクレームを持つ SAML アサーションが含まれる
    * クレームとは、ユーザーとそのグループに関する情報で、クレームルールはその情報を SAML 属性にマッピングする
    * これにより、AWS が IAM ポリシー内でフェデレーティッドユーザーのアクセス許可を確認するのに必要な属性が、IdP からの SAML 認証レスポンスに確実に含まれる

### 認証レスポンスの SAML アサーションを設定
* 組織内でユーザーの ID が確認されたら、外部 ID プロバイダー (IdP) は AWS SAML エンドポイント (https://signin.aws.amazon.com/saml) に認証レスポンスを送信
  * このレスポンスは、HTTP POST Binding for SAML 2.0 標準に従った SAML トークンを含み、さらに以下の要素またはクレームを含む POST リクエストであることが必要
  * 設定方法は、やはり IdP 次第
* [認証レスポンスの SAML アサーションを設定する](https://docs.aws.amazon.com/ja_jp/IAM/latest/UserGuide/id_roles_providers_create_saml_assertions.html#saml-attribute-mapping)

## AWS CLI
* SSO でログイン
  * SAML Response を取得しておく（```response.log``` として）

``` bash
$ aws sts assume-role-with-saml --role-arn arn:aws:iam::<account_id>:role/<role_name> --principal-arn arn:aws:iam::<account_id>:saml-provider/<provider> --saml-assertion file://response.log --duration-seconds 43200
````

* ~/.aws/credentials
  * 上の結果を設定
```
[default]
aws_access_key_id=...
aws_secret_access_key=...
aws_session_token=...
```
* ~/.aws/config
  * ReadOnly は適当に用意しておいた Role
```
[profile default]

[profile readonly]
role_arn = arn:aws:iam::<account_id>:role/ReadOnly
source_profile = default
```

* 動作確認
``` bash
# それぞれちゃんとうまく動いていることがわかる
$ aws sts get-caller-identity
$ aws sts get-caller-identity --profile readonly
# default だとできなくて、readonly だとできる操作
$ aws s3 ls --profile readonly
```

* これを設定すると、aws-sdk-go とか AWS SDK を利用したツールでも対応可能
```
export AWS_PROFILE=readonly
export AWS_SDK_LOAD_CONFIG=true
```

## Reference
* [SAML](https://techinfoofmicrosofttech.osscons.jp/index.php?SAML)
* [IAM SAML ID プロバイダーの作成](https://docs.aws.amazon.com/ja_jp/IAM/latest/UserGuide/id_roles_providers_create_saml.html)
* [Keycloakを使ってAWSにSSO接続する方法](https://www.ryuzee.com/contents/blog/7129a)