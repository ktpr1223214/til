---
title: Java
---

## Java
### Java どれ選べば？
* [Java有償化（していない件）について](https://speakerdeck.com/gishi_yama/java-do-number-osc19do)
* [Javaは今でも無償です、という話 / Java is still free](https://speakerdeck.com/kishida/java-is-still-free)

* adoptopenjdk が良さげ？

### 環境構築
* adoptopenjdk
    * ver は11
``` bash
# install
# そのままでは（20190614時点）12が入る
$ brew tap AdoptOpenJDK/openjdk 
$ brew cask install adoptopenjdk11
$ java -version

# uninstall
# JDK パス確認
$ /usr/libexec/java_home -V
$ sudo rm -rf /Library/Java/JavaVirtualMachines/<some jdk>.jdk
```

* intellij
    * SDK のパスを指定
    * /usr/libexec/java_home -V の結果の、~.jdk フォルダを指定すれば読み込んでくれる    

* gradle
``` bash
$ brew install gradle
$ gradle -v
```

* spring


### Gradle
* [Gradle使い方メモ](https://qiita.com/opengl-8080/items/4c1aa85b4737bd362d9e)
``` bash
# build/ に生成
$ gradle build
# build/ 削除
$ gradle claen
```

``` groovy
apply plugin: 'java'
```
これを指定し、gradle build した後にクラスを実行するには、
``` bash
# クラスパスを明示的に指定して実行
$ java -cp build/classes/java/main/ hello.HelloWorld
```

``` groovy
# gradle から実行(gradle run)できるようにするのに必要 
apply plugin: 'application'

mainClassName = 'hello.HelloWorld'

# javac の -source/-target に対応
# -source: コンパイルするソースが準拠している JDK のバージョン
# -target: コンパイルして出来た class が実行できる JavaVM JRE のバージョンを指定
sourceCompatibility = 1.8
targetCompatibility = 1.8

# 
```


### checkstyle
* lint
    * checkstyle
*     
 

## Reference
   
