---
title: Secret Manager
---

## 権限
* [Access control](https://cloud.google.com/secret-manager/docs/access-control)
  * ここのページの Note にあるように、roles/owner は secretmanager.versions.access を持つが、roles/editor と roles/viewer は持たない
  * 特に roles/editor が持たないので、Cloud Functions がデフォルトでは secertmanager のバージョンにアクセスできないことに注意
