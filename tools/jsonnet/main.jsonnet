local base = import "base.libsonnet";

base.new2(title="hoge")
// こうすると JSON つながる
{
    a: false
}
+
{
    b: true
}
// 明示的につなげることも勿論可能
+
{
    b: false
}
// 後の方で上書き
+
{
    x: base.id("x"),
    y: base.x
}
+
base.new3().method("hoge")
+
base.hoge.add()
// base.hoge
