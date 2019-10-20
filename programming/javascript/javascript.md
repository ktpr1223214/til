---
title: Javascript
---

## Javascript

### setup(macos)
* anyenv で入れる
    * anyenv update も入れる

``` bash
$ anyenv install ndenv
$ exec $SHELL -l
# install 可能一覧
$ ndenv install -l
$ ndenv install v.~.~

# 確認
$ ndenv versions
# global 設定
$ ndenv global v.~.~
# local 設定
$ ndenv local v.~.~
```

``` bash
# yarn yarn は brew で
$ brew install yarn --ignore-dependencies
# upgrade したいとき
# brew upgrade yarn すると、結局 node をいれてしまうので、uninstall から再度 install が手っ取り早い？(もう少しましな方法ありそうだけど)
$ brew uninstall yarn
$ brew uninstall yarn --force yarn
```

### project setup
* yarn を使う
``` bash
# init
$ yarn init -y

# webpack
$ yarn add webpack webpack-cli --dev
$ yarn add url-loader file-loader css-loader style-loader --dev
$ yarn add typings-for-css-modules-loader --dev
$ yarn add mini-css-extract-plugin --dev
# for minifiy css/js
$ yarn add optimize-css-assets-webpack-plugin --dev
$ yarn add uglifyjs-webpack-plugin --dev

# webpack-dev-server
$ yarn add webpack-dev-server --dev
```
