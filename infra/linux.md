---
title: Linux
---

## user/group
* ユーザーの情報は /etc/passwd
  * ユーザー名:パスワード:ユーザーID:グループID:その他の情報:ホームディレクトリ:シェル
* グループの情報は /etc/group
  * グループ名:パスワード:グループID:ユーザーリスト
* ```man passwd``` や ```man group``` で詳細確認

## wheel/sudo
* wikipedia より
The wheel group is a special user group used on some Unix systems, mostly BSD systems, to control access to the su or sudo command, which allows a user to masquerade as another user (usually the super user). Debian-like operating systems create a group called sudo with similar purpose to wheel group.

* visudo
  * 書式: ユーザー ホスト=(権限) コマンド
    * NOPASSWD and PASSWD: パスワード必要 or not
    * NOEXEC and EXEC: sudo が noexec サポートつきでコンパイルされ、 使用しているオペレーティングシステムがそれに対応している場合、NOEXEC タグを利用すれば、 動的にリンクされた実行ファイルが、そこからさらにコマンドを実行するのを防ぐ
      * ある Cmnd にタグをセットすると、 Cmnd_Spec_List 中のそれ以後の Cmnd は、 反対の意味を持つタグによって変更されないかぎり、そのタグを継承することになる (すなわち、PASSWD は NOPASSWD を無効にし、NOEXEC は EXEC を無効に
* Amazon Linux 2 では ```/etc/sudoers.d/90-cloud-init-users``` を更新
  * ```sudo visudo -f /etc/sudoers.d/90-cloud-init-users```

## /etc/security/limits.conf
* ユーザーごとのリソースを制限できるファイル
* 記述例
  * ```*``` と別に root は指定が必要な模様
  * https://bugs.launchpad.net/ubuntu/+source/pam/+bug/65244
    * 下の man コマンドでは該当するような記述がないが、ここによればドキュメント不備らしい？
```
root            hard    nofile      65536
root            soft    nofile      65536
*               hard    nofile      65536
*               soft    nofile      65536
```
* ```man limits.conf``` から
  * 設定可能項目
    * cpu: maximum CPU time (minutes)
    * nproc: maximum number of processes
    * as: address space limit (KB)
    * maxlogins: maximum number of logins for this user (this limit does not apply to user with uid=0)
    * maxsyslogins: maximum number of all logins on system; user is not allowed to log-in if total number of all users' logins is greater than specified number (this limit does not apply to user with uid=0)
    * priority: the priority to run user process with (negative values boost process priority)
    * locks: maximum locked files (Linux 2.4 and higher)
    * sigpending: maximum number of pending signals (Linux 2.6 and higher)
    * msgqueue: maximum memory used by POSIX message queues (bytes) (Linux 2.6 and higher)
    * nice: maximum nice priority allowed to raise to (Linux 2.6.12 and higher) values: [-20,19]
    * rtprio: maximum realtime priority allowed for non-privileged processes (Linux 2.6.12 and higher)
  * soft/hard:
    * hard: for enforcing hard resource limits. These limits are set by the superuser and enforced by the Kernel. The user cannot raise his requirement of system resources above such values
    * soft: for enforcing soft resource limits. These limits are ones that the user can move up or down within the permitted range by any pre-existing hard limits. The values specified with this token can be thought of as default values, for normal system usage

## MTA（Mail Transfer Agent）
* CentOS では、sendmail/postfix をいつでも切り替えられるらしい
  * alternatives という仕組みで実現している
``` bash
# 現在選択している MTA の表示
$ alternatives --display mta

# postfix に MTA 切り替え（alternatives --config mta で、対話的に変更も可能だが、ワンライナーだと以下のように）
$ alternatives --set mta /usr/sbin/sendmail.postfix
```

## /etc/aliases
* sendmail 用のエイリアスを定義するファイル
  * root: hoge という書式で、root 宛のメールが hoge に回送される
  * サーバドメインが、hoge.co.jp ならば root@hoge.co.jp へのメールが hoge@hoge.co.jp へ回送されるということ

``` bash
# /etc/aliases を編集した場合には、このコマンド実行が必要
$ newaliases
```

## /etc/skel
* /etc/skel ディレクトリは、新規ユーザーの home directory にコピーされる
  * useradd コマンドによる作成時に
  * /etc/defualt/useradd で設定される

## Reference
* [GCC Inline Assembler](http://caspar.hazymoon.jp/OpenBSD/annex/gcc_inline_asm.html)
  * Kernel のコードを読むときの前提知識として
* [A Heavily Commented Linux Kernel Source Code](http://oldlinux.org/download/ECLK-5.0-WithCover.pdf)
* [AT&T Assembly Syntax](http://web.archive.org/web/20080215230650/http://sig9.com/articles/att-syntax)
  * 'quick-n-dirty' introduction to the AT&T assembly language syntax