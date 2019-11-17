local a = 2;

// local で function を定義しても、library 的には使えないはずなので
// こういう :: を使った定義となる
// また :: で hidden にしないと function は JSON に出てこれない
{
  new(title)::
    {
      title: title + a,
    },
  id:: function(x) x,
  x:: 2,
}

// 複数もいける
{
  new2(title)::
    {
      title: title,
    },
}

{
  new3()::
    {
       method(title):: {
         target: title
       }
    }
}

{
    hoge:: {
        _next: 0,
        add():: self {
            local x = super._next,
            _next: x + 1,
            char: std.char(std.codepoint('A') + x),
        },
    },
}