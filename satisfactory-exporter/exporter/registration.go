package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricVectorDetails struct {
	Name   string
	Help   string
	Labels []string
}

var RegisteredMetricVectors = []MetricVectorDetails{}
var RegisteredMetrics = []*prometheus.GaugeVec{}

func RegisterNewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.Desc {
	return prometheus.NewDesc(
		opts.Name,
		opts.Help,
		labelNames,
		nil,
	)
}
