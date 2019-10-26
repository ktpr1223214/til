---
title: Grafana
---

## Grafana JSON model
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

## Reference
* [Dashboard JSON](https://grafana.com/docs/reference/dashboard/#dashboard-json)
    * JSON Model について
* [Jsonnet Tutorial](http://35.190.68.35/docs/tutorial.html)
    * <function>::... について書いてある