---
title: JSON
---

## JSON
* [RFC7159](https://tools.ietf.org/html/rfc7159)
* [Introducing JSON](http://json.org/)
    * A value can be a string in double quotes, or a number, or true or false or null, or an object or an array
    * なので、"ok" は正しい形式だが、ok は駄目(text/html とかで返ってきているのを1回勘違いしてたことがあるので、注意)
    * 数字はそのままで良く、true/false/null も問題ない
    
* Go で stream に対する decode 
``` go
var out interface{}
err := json.NewDecoder(strings.NewReader("truetrue")).Decode(&out)
fmt.Println(out)
// true

var out interface{}
err := json.NewDecoder(strings.NewReader("nullll")).Decode(&out)
fmt.Println(out)
// nil
```