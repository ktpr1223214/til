---
title: Nil
---

## Nil
### 概要
* Variables of interface type also have a distinct dynamic type, which is the concrete type of the value assigned to the variable at run time (unless the value is the predeclared identifier nil, which has no type)
    * interface 型の変数は動的に具体的な型をもつが(runtime に)、nil の場合は例外(nil has no type)
* nil is a predeclared identifier representing the zero value for a pointer, channel, func, interface, map, or slice type.
    * [nil](https://golang.org/pkg/builtin/#pkg-variables)
    * zero 「value」
* interface 以外はそこまで悩むところはない。interface は (type, value) の組からなり、どちらも(nil, nil) の場合のときのみ nil
* The typed nil emerges as a result of the definition of a interface type; a structure which contains the concrete type of the value stored in the interface, and the value itself

## Reference
* []()