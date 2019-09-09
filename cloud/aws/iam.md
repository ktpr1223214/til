---
title: IAM
---

## IAM

### 用語
* IAMユーザー
    * UserName
    * Path
* IAMグループ
    * GroupName
    * Path
* IAMロール
    * RoleName
    * Path
    * AssumeRolePolicyDocument

* IAMプリンシパル
    * リソースへのアクセスが許可または拒否されるエンティティとして定義
* アイデンティティベース（IDベース）のポリシー
    * インラインポリシーと管理ポリシー
* リソースベースのポリシー
    * 操作を行われるモノ(AWSリソース)に関連付けるポリシー
    * サポートされるリソースにインラインポリシーをアタッチ
    * 実行の主体として、Principalを指定する必要がある
    * 例えば、S3はこのポリシーを付与できるがDynamoDBではできない
* 信頼ポリシー
    * ロールを引き受けるユーザーを定義する JSON 形式のドキュメント

### インスタンスプロファイル
* インスタンスプロファイルは IAM ロールのコンテナであり、インスタンスの起動時に EC2 インスタンスにロール情報を渡すために使用
    * terraform などを使わないと EC2 インスタンスと IAM role が直接紐づくと考えてしまうが、実際はインスタンスプロファイルを介している
* [EC2にIAMRole情報を渡すインスタンスプロファイルを知っていますか？](https://dev.classmethod.jp/cloud/aws/do_you_know_iaminstanceprofile/)      

## terraform

* IAM user(inline policy)
``` terraform
# IAMユーザーにinline policyを紐付ける
resource "aws_iam_user_policy" "lb_ro" {
  name = "test"
  user = "${aws_iam_user.lb.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# user
resource "aws_iam_user" "lb" {
  name = "loadbalancer"
  path = "/system/"
}

# access key
resource "aws_iam_access_key" "lb" {
  user = "${aws_iam_user.lb.name}"
}
```

* IAM user(managed policy)

``` terraform
resource "aws_iam_user" "admin" {
  name = "admin"
}

# IAMユーザーに managed policy を紐付ける
resource "aws_iam_user_policy_attachment" "admin" {
  user = "${aws_iam_user.admin.name}"
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
```

* 信頼ポリシー

``` hcl-terraform
# 上2つがAssumeRolePolicyDocumentに相当
data "aws_iam_policy_document" "instance-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

# assume_role_policy :=  The policy that grants an entity permission to assume the role.
resource "aws_iam_role" "instance" {
  name               = "instance_role"
  path               = "/system/"
  assume_role_policy = "${data.aws_iam_policy_document.instance-assume-role-policy.json}"
}

# roleに権限設定
resource "aws_iam_role_policy" "instance-policy" {
    name = "instance-policy"
    role = "${aws_iam_role.instance.id}"

    policy = "${data.aws_iam_policy_document.instance-policy.json}"
}

// json形式のIAM policy
data "aws_iam_policy_document" "instance-policy" {
    statement {
        sid = "AllowAccessToCloudWatchLogs"
        effect = "Allow"
        actions = [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
        ]
        resources = ["*"]
    }
}
```