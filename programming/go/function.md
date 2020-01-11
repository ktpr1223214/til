---
title: function
---

## function

### 引数
* [Be wary of functions which take several parameters of the same type](https://dave.cheney.net/2019/09/24/be-wary-of-functions-which-take-several-parameters-of-the-same-type)
* 同じ型の引数を複数取る関数の API は誤って使いやすいので、その対策
``` golang
// 以下2つは誤りやすい
CopyFile("/tmp/backup", "presentation.md")
CopyFile("presentation.md", "/tmp/backup")

// ので、helper type を導入すると吉という話
type Source string

func (src Source) CopyTo(dest string) error {
	return CopyFile(dest, string(src))
}

func main() {
	var from Source = "presentation.md"
	from.CopyTo("/tmp/backup")
}
```
