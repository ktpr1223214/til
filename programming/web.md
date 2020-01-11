* editが成功したら詳細画面へリダイレクトするのが主流
* ほかのcontrollerのメソッドを見てもらうといいんですが、まず最初に is_valid_post() をチェックして、そうでなかったらGETと同じ処理にする、というのが主流です。
  現状のコードだと、 is_valid_post ではないPOSTをすると edit() が何もreturnしなくて500エラーになりますね

## Session
* サーバ側でもログイン状態を持つという観点もある
    * サーバ側で消せば、クッキーがあろうとログインはできない

* cookie がセットされていない場合
    * http.StatusUnauthorized

## Cookie

## Endpoint
* [Monitoring the health of your application - The upgraded "/ping" route ](https://www.sohamkamani.com/blog/architecture/2018-09-06-application-health-monitoring/)