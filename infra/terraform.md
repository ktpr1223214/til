---
title: Terraform
---

## terraform

## tfstate
* [Terraform state 概論](https://chroju.github.io/blog/2019/12/13/terraform_state_introduction/)を参考に
  * tfstate からは、先述のように各 Terraform resource と対応する現実のリソースが存在するのかどうかを読み取るだけ
  * terraform plan のパターン
    * tffile に: ある tfstate に: ある 現実のリソースと tffile の差異が: ない -> plan 結果は: No changes
    * tffile に: ある tfstate に: ある 現実のリソースと tffile の差異が: ある -> plan 結果は: change
    * tffile に: ある tfstate に: ない 現実のリソースと tffile の差異が: 差異は確認しない -> plan 結果は: add
    * tffile に: ない tfstate に: ある 現実のリソースと tffile の差異が: 差異は確認しない -> plan 結果は: destroy
