---
title: GCP
---

## gcloud
``` bash
# auth
$ gcloud auth login

# デフォルトのプロジェクト設定
$ gcloud config set project <project-id>

# デフォルトのリージョン設定
# 確認は gcloud compute zones list
$ gcloud config set compute/region <region>

# デフォルトのコンピューティングゾーン設定
# ex. asia-northeast1-a
$ gcloud config set compute/zone <compute-zone>

# 設定確認
$ gcloud projects list
$ gcloud config list
```
