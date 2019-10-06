---
title: Shell
---

## basic
``` bash
# コマンド1が成功した場合のみ、コマンド2が実行
$ command1 && command2
# コマンド1が失敗した場合のみ、コマンド2が実行
$ command1 || command2
```

## grep
``` bash
# -e は
$ grep -l "文字列" *.go | xargs sed -i -e 's/変換前文字列/変換後文字列/g'
```

## EOF
``` bash
# 標準出力に
cat <<EOF
{
    "uid": "uin",
    "title": "title"
}
EOF

# ファイルを作りたければリダイレクト
cat <<EOF > test.txt
{
    "uid": "uin",
    "title": "title"
}
EOF
```
