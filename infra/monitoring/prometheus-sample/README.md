# prometheus-sample

## 挙動確認
``` bash
# まずは webhook を受けるサンプルを起動
$ go run cmd/webhook/main.go

# まずは単一の critical 発火
cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/test/service/fuga
# TYPE up counter
up 0
EOF

# ここで待てば、repeat_interval の挙動検証が可能(repeat_interval は適当に短くしないといけないだろうが)

# 次に、同一 service の critical なアラートを発火
# これは纏められる
cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/node_exporter
# TYPE up counter
up 0
EOF

# 次に、異なる service の critical なアラート発火
# これは纏められず別のアラートグループとなる
cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/test/service/hoge
# TYPE up counter
up 0
EOF

# これは、2つの別 ticket 発火する
# なぜなら、alertname での group_by も追加しているから
cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/test/service/piyo
# TYPE up counter
up 0
EOF
```
