---
title: RDS
---

## RDS

## subnet group
* https://aws.amazon.com/jp/premiumsupport/knowledge-center/rds-db-subnet-group/

## serverless
``` bash
# version の指定がややこいので、現在使える組合せを調べて指定
$ aws rds describe-db-engine-versions --engine aurora --query "DBEngineVersions[].SupportedEngineModes"
```
