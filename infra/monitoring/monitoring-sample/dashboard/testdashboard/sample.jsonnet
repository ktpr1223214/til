local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local template = grafana.template;
local singlestat = grafana.singlestat;
local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

local nodeMetrics = import 'node_exporter.libsonnet';
local basic = import 'base_graph.libsonnet';

local buildInfo =
  singlestat.new(
    title='Version',
    datasource='Prometheus',
    format='none',
    valueName='name',
  ).addTarget(
    prometheus.target(
      "prometheus_build_info{instance='$instance'}",
      legendFormat='{{ version }}',
    )
  );

local systemLoad =
  singlestat.new(
    title='5m system load',
    datasource='Prometheus',
    format='none',
    valueName='current',
    decimals=2,
    sparklineShow=true,
    colorValue=true,
    thresholds='4,6',
  ).addTarget(
    prometheus.target(
      "node_load5{instance='$instance'}",
    )
  );

local networkTraffic =
  graphPanel.new(
    title='Network traffic on eth0',
    datasource='Prometheus',
    linewidth=2,
    format='Bps',
    aliasColors={
      Rx: 'light-green',
      Tx: 'light-red',
    },
  ).addTarget(
    prometheus.target(
      "rate(node_network_receive_bytes_total{instance='$instance',device='eth0'}[1m])",
      legendFormat='Rx',
    )
  ).addTarget(
    prometheus.target(
      "irate(node_network_transmit_bytes_total{instance='$instance',device='eth0'}[1m]) * (-1)",
      legendFormat='Tx',
    )
  );

local samplets =
  basic.timeseries(
    title='Disk Write Total Time',
    description='Total time spent in write operations across all disks on the node. Lower is better.',
    // fqdn は取れている前提？なので、instance とか？
    query='sum(rate(node_disk_write_time_seconds_total{instance="$instance"}[$__interval])) by (fqdn)',
    legendFormat='{{ instance }}',
    format='s',
    interval='30s',
    intervalFactor=1,
    yAxisLabel='Total Time/s',
    legend_show=false,
    linewidth=1
  );


dashboard.new(
  'Prometheus test',
  tags=['prometheus'],
  schemaVersion=18,
  editable=true,
  time_from='now-1h',
  refresh='1m',
)

# https://github.com/grafana/grafonnet-lib/blob/69bc267211790a1c3f4ea6e6211f3e8ffe22f987/grafonnet/template.libsonnet#L69
# PROMETHEUS_DS という名前の Type: DataSource の variables を
# prometheus DataSource
.addTemplate(
  template.datasource(
    'PROMETHEUS_DS',
    'prometheus',
    'Prometheus',
    hide='label',
  )
)

# https://github.com/grafana/grafonnet-lib/blob/69bc267211790a1c3f4ea6e6211f3e8ffe22f987/grafonnet/template.libsonnet#L2
# instance という名前の variables を $PROMETHEUS_DS という datasource(上で定義した)から、
# label_values(prometheus_build_info, instance) というクエリで定義
# label は、The name of the dropdown for this variable.
# これは、variable types としては、Query にあたる
.addTemplate(
  template.new(
    'instance',
    '$PROMETHEUS_DS',
    // 'label_values(prometheus_build_info, instance)',
    'label_values(instance)',
    label='Instance',
    refresh='time',
  )
)
.addPanel(nodeMetrics.nodeMetricsDetailRow("instance='$instance'"), gridPos={ x: 0, y: 6000 })
.addPanels(
  [
    // https://grafana.com/docs/reference/dashboard/#panel-size-position
    // 書き方としては、buildInfo の gridPos 要素を上書きと理解すれば ok ?
    buildInfo { gridPos: { h: 4, w: 3, x: 0, y: 0 } },

    systemLoad { gridPos: { h: 4, w: 4, x: 3, y: 0 } },

    networkTraffic { gridPos: { h: 8, w: 7, x: 0, y: 4 } },

    samplets { gridPos: { h: 8, w: 7, x: 0, y: 13 } },
  ]
)
