---
title: sudoers
---

## sudoers
* <user or group> ALL = (ALL:ALL) ALL
  * who where = (as_whom) what
    * 誰が どのホストで = (誰になって、どのグループで) 何を
      * () の二つ目の : 以降はグループを指定
  * <user or group>: ユーザー名を指定、もしくは % 付きで group 名
* 例
  * ec2-user ALL=(ALL) NOPASSWD:ALL
    * NOPASSWD はパスワードの確認が不要
    * ec2-user は、どこからでもすべてのユーザーとしてすべてのコマンドを実行可能
  * piyota ALL=(akuma) ALL
    * piyota は akuma として全てのコマンドを実行可能
    * sudo -u akuma ls はできて、sudo ls は出来ない（= root ユーザとしての実行）
* /etc/sudoers 配下を編集する時は visudo を使うこと

## Reference
* [/etc/sudoersとは](https://wa3.i-3-i.info/word13805.html)
* [Sudoers Manual](https://www.sudo.ws/man/1.8.18/sudoers.man.html)