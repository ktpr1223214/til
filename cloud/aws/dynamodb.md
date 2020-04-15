---
title: DynamoDB
---

## tips
* [list_append && if_not_exists](https://stackoverflow.com/questions/34951043/is-it-possible-to-combine-if-not-exists-and-list-append-in-update-item)
* aws-sdk-go での expression の確認
``` go
expr, err := expression.NewBuilder().WithUpdate(
	expression.Set(
		expression.Name("hoge"),
		expression.ListAppend(expression.Name("hoge"),
			expression.Value([]string{something}))),
).Build()
if err != nil {
	return err
}
// この辺りをみると、条件式の中身がわかる
log.Println(*expr.Update())
log.Println(expr.Names())
log.Println(expr.Values())
```

## aws-sdk-go 関連
* [aws-sdk-go での empty list](https://msanatan.com/2018/08/31/dynamodb-lambdas-go-and-an-empty-list/)
