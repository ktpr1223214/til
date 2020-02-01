---
title: Grafana
---

## Grafana JSON model
* links
    * graph にリンクを付けられる
``` json
"links": [
    {
        "title": "hoge",
        "url": "http://hoge.com"
    }
]
```    

* targets
    * Query を書くところ
```` json
"targets": [
    {
      "expr": "prometheus_build_info{instance=\"$instance\"}",
      "format": "time_series",
      "intervalFactor": 2,
      "legendFormat": "{{ version }}",
      "refId": "A"
    }
]
````    

### templating
* 設定
    * General
        * Name: Prometheus_DS
        * Type: Datasource
    * Data source options
        * Type: Prometheus         

``` json
 "templating": {
    "list": [
      {
        "current": {
          "tags": [],
          "text": "Prometheus",
          "value": "Prometheus"
        },
        "hide": 0,
        "includeAll": false,
        "label": null,
        "multi": false,
        "name": "Prometheus_DS",
        "options": [],
        "query": "prometheus",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "type": "datasource"
      }
    ]
  },
```

## grafonnet
### prometheus.libsonnet
```
// target は、query を書く
{
  target(
    expr,
    format='time_series',
    intervalFactor=2,
    legendFormat='',
    datasource=null,
    interval=null,
    instant=null,
    hide=null,
  ):: {
    [if hide != null then 'hide']: hide,
    [if datasource != null then 'datasource']: datasource,
    expr: expr,
    format: format,
    intervalFactor: intervalFactor,
    legendFormat: legendFormat,
    [if interval != null then 'interval']: interval,
    [if instant != null then 'instant']: instant,
  },
}
```

### graph_panel.libsonnet
* 基本の graph panel
* 細かな設定は以下ページ参照 
    * https://grafana.com/docs/features/panels/graph/
    
``` 
addTarget(target):: self {
  // automatically ref id in added targets.
  // https://github.com/kausalco/public/blob/master/klumps/grafana.libsonnet
  local nextTarget = super._nextTarget,
  _nextTarget: nextTarget + 1,
  // trget は A~ の ID を振られていることから(画面みたらわかる)
  targets+: [target { refId: std.char(std.codepoint('A') + nextTarget) }],
},
```    
    

### link.libsonnet
* dashboard 上部の link(graph に紐づくやつではない)
    * type: "dashboards"/"link"
        * 前者はダッシュボード一覧が生える
    * title/url は type: "link" で使う
    * icon: external link とか dashboard          
``` 
{
  dashboards(
    title,
    tags,
    asDropdown=true,
    includeVars=false,
    keepTime=false,
    icon='external link',
    url='',
    targetBlank=false,
    type='dashboards',
  )::
    {
      asDropdown: asDropdown,
      icon: icon,
      includeVars: includeVars,
      keepTime: keepTime,
      tags: tags,
      title: title,
      type: type,
      url: url,
      targetBlank: targetBlank,
    },
}
```

### row.libsonnet
    * panels をまとめる単位
``` 
{
  new(
    title='Dashboard Row',
    height=null,
    collapse=false,
    repeat=null,
    showTitle=null,
    titleSize='h6'
  ):: {
    collapse: collapse,
    collapsed: collapse,
    [if height != null then 'height']: height,
    panels: [],
    repeat: repeat,
    repeatIteration: null,
    repeatRowId: null,
    showTitle:
      if showTitle != null then
        showTitle
      else
        title != 'Dashboard Row',
    title: title,
    type: 'row',
    titleSize: titleSize,
    addPanels(panels):: self {
      panels+: panels,
    },
    addPanel(panel, gridPos={}):: self {
      panels+: [panel { gridPos: gridPos }],
    },
  },
}
```

### templates.libsonnet

```
// template.datasource の引数は上から、name/query/current
// また、type(Defines the variable type.)は、"datasource" 固定
// なので、PROMETHEUS_DB という名前で、 
ds:: template.datasource(
    'PROMETHEUS_DS',
    'prometheus',
    'Prometheus',
    regex='/(.*-gprd|Global)/',
)

// template.new の引数は上から、name/datasource/query
// type は "query" 固定
// This variable type allows you to write a data source query that usually returns a list of metric names, tag values or keys. For example, a query that returns a list of server names, sensor ids or data centers.

```

### saturation
```
// clamp_max(v instant-vector, max scalar) clamps the sample values of all elements in v to have an upper limit of max 
'clamp_min(clamp_max(' + query + ',1),0)' 
```

### stacking && null value 
Stack - Each series is stacked on top of another
Percent - Available when Stack are checked. Each series is drawn as a percentage of the total of all series
Null value - How null values are displayed

### addTarget 関数
* 各種グラフ毎に、関数として実装されているはず
    * cf. https://github.com/grafana/grafonnet-lib/blob/f471d13189822c9cf94d5f198fdced3323279877/grafonnet/singlestat.libsonnet#L124 
* こういうふうに、new の返りの中で hidden で定義されているので、使わなければ消すし、使う場合はメソッドチェイン的に呼び出せるのかと思う

### $__interval
* https://grafana.com/docs/reference/templating/#interval

## examples
* https://dashboards.gitlab.com/d/web-main/web-overview?orgId=1
    * https://gitlab.com/gitlab-com/runbooks/blob/master/dashboards/web/main.dashboard.jsonnet    

## Reference
* [Dashboard JSON](https://grafana.com/docs/reference/dashboard/#dashboard-json)
    * JSON Model について
* [Jsonnet Tutorial](http://35.190.68.35/docs/tutorial.html)
    * <function>::... について書いてある
* [slo-libsonnet](https://github.com/metalmatze/slo-libsonnet)
    * SLO と libsonnet    