---
title: API
---

## API

## List のページ分割
* 後からサポートすると、クライアントとの連携が面倒なので結果が小さくてもサポートしておくほうが安全
* List でページ分割の実装(cf. [リストのページ分割](https://cloud.google.com/apis/design/design_patterns#list_pagination))
  * List メソッドのリクエスト メッセージで string フィールド page_token を定義します。クライアントはこのフィールドを使用して、リスト結果の特定のページをリクエストします。
  * List メソッドのリクエスト メッセージで int32 フィールド page_size を定義します。クライアントはこのフィールドを使用して、サーバーによって返される結果の最大数を指定します。1 ページ内に返される結果の最大数が、サーバーによってさらに制限されることがあります。page_size が 0 の場合、返される結果の数はサーバーが決定します。
  * List メソッドのレスポンス メッセージで string フィールド next_page_token を定義します。このフィールドは、結果の次のページを取得するためのページ設定トークンを表します。値が "" の場合、リクエストに対してそれ以上の結果がないことを意å³します。

### pagination の参考
* [Pagination](https://www.django-rest-framework.org/api-guide/pagination/)
  * DRF のページだが、一般に良さげ
* [Building Cursors for the Disqus API](https://cra.mr/2011/03/08/building-cursors-for-the-disqus-api)
* [Evolving API Pagination at Slack](https://slack.engineering/evolving-api-pagination-at-slack-1c1f644f8e12)
  * まとまっていてわかりやすい

## Health check
* [Health Check Response Format for HTTP APIs](https://tools.ietf.org/html/draft-inadarei-api-health-check-03)

## Reference
* [API 設計ガイド](https://cloud.google.com/apis/design/)
