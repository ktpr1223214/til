---
title: Interface
---

## Interface

## Reference
* [How much interface do I need?](https://blog.chewxy.com/2018/03/18/golang-interfaces/)
    * Define the types
    * Define the interface at point of use.
        * これらによって、依存が減らせる → より robust になる
    * func funcName(a INTERFACETYPE) CONCRETETYPE
    * また、先立って interface を定義しても良い場合として以下を挙げている(例外の無い規則は無い)
        * Sealed interfaces
        * Abstract data types
        * Recursive interfaces

* [Preemptive Interface Anti-Pattern in Go](https://medium.com/@cep21/preemptive-interface-anti-pattern-in-go-54c18ac0668a)
    * Go/Java の interface の大きな違いとして、implicit/explicit が挙げられる
    * 先立っての interface 定義が Go ではあまり良い方法ではない(カーニハン本でも、interface の末に書いてある話)ことの説明として分かりやすい
    * accept interfaces, return structs だが、accept structs, return structs でも Go なら後から変更をするのは容易
        * なぜなら、使う側の引数を interface にするだけで良く、元実装での対応が必要ないから
            * 勿論常にそうとは限らないとは思うが..(interface で使い分けるために、複数実装の（文字通りの）インターフェースを変更したりといった修正はあるかも)            

* [What “accept interfaces, return structs” means in Go](https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8)
    * Remove unneeded abstractions
        * ここが？
    * Ambiguity of a user’s need on function input
        * This imbalance between being able to precisely control the output, but be unable to anticipate the user’s input, creates a stronger bias for abstraction on the input than it does on the output.
    * Simplify function inputs
        * 構造体の他のメソッドは捨象出来る
    * interface を一般的にするというのは、philosophy の本とも合致か

* [How much interface do I need?](https://peter.bourgon.org/go-for-industrial-programming/#how-much-interface-do-i-need)
    * interfaces are consumer contracts, not producer (implementor) contracts—so, as a rule, we should be defining them at callsites in consuming code, rather than in the packages that provide implementations

* [Does a concrete type implement an interface in Go?](https://eli.thegreenplace.net/2019/does-a-concrete-type-implement-an-interface-in-go/)
    * interface の compile-time/run-time check 話を軽く
        ``` go var _ Munger = (*Foo)(nil)```

* [Go and Algebraic Data Types](https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/)
    * Go に Algebraic data types は存在しないが、(elegant でないにしろ)同じようなことは出来るという話

### implementation など
* [Go Data Structures: Interfaces](https://research.swtch.com/interfaces)
    * russ cox
* 