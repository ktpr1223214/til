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
# -r つけると、recursive(なので、grep -rl とかにすれば)
$ grep -l "文字列" *.go | xargs sed -i 'sample.bak' 's/変換前文字列/変換後文字列/g'
# macos で、上書きする場合(-i オプションの挙動が BSD のため異なることに注意)
$ grep -rl "'" ./common/ | xargs sed -i '' 's/変換前文字列/変換後文字列/g'
# ex. ' を " に修正
$ grep -rl "'" ./common/ | xargs sed -i '' "s/'/\"/g"
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
