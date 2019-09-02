---
title: Go Test
---

## Go Test

## tips
* [from @peterbourgon tweet](https://twitter.com/peterbourgon/status/1151886045861928961)
``` 
Underutilized testing strategy: test your log output. This forces you to do several things you should be doing anyway:

- Loggers are dependencies, not globals
- Routable log output (to in-mem buffers)
- Well-defined log schemas (structured logging)
- Deterministic components

Testing log output can also be a good proxy for asserting internal state/behavior, 
without having to get your unit test dirty by poking at private data.

It’s worth clarifying that testing log output should be, at most, just part of your overall testing breakfast.

@benbjohnson
Another handy tip is to check `testing.Verbose()` and use `io.MultiWriter` to send logs to `os.Stderr` and your in-memory buffer so you can debug logs with the `go test -v` flag.
```

## Reference
* [Testing Your (HTTP) Handlers in Go](https://blog.questionable.services/article/testing-http-handlers-go/)
* [](https://deliveroo.engineering/2019/05/17/testing-go-services-using-interfaces.html#services)
