package exporter

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type ProductionCollector struct {
	frmTarget string
	logger    log.Logger
}

type ProductionDetails struct {
	ItemName           string  `json:"Name"`
	ProdPercent        float64 `json:"ProdPercent"`
	ConsPercent        float64 `json:"ConsPercent"`
	CurrentProduction  float64 `json:"CurrentProd"`
	CurrentConsumption float64 `json:"CurrentConsumed"`
	MaxProd            float64 `json:"MaxProd"`
	MaxConsumed        float64 `json:"MaxConsumed"`
}

func NewProductionCollector(frmApiAddress string, logger log.Logger) *ProductionCollector {
	return &ProductionCollector{
		frmTarget: frmApiAddress + "/getProdStats",
		logger:    logger,
	}
}

func (c ProductionCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *ProductionCollector) Collect(ch chan<- prometheus.Metric) {
	details := []ProductionDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading production statistics from Ficsit Metrics", "err", err)
		return
	}

	for _, d := range details {
		ch <- prometheus.MustNewConstMetric(ItemsProducedPerMin, prometheus.GaugeValue, d.CurrentProduction, d.ItemName)
		ch <- prometheus.MustNewConstMetric(ItemsConsumedPerMin, prometheus.GaugeValue, d.CurrentConsumption, d.ItemName)
		ch <- prometheus.MustNewConstMetric(ItemProductionCapacityPercent, prometheus.GaugeValue, d.ProdPercent, d.ItemName)
		ch <- prometheus.MustNewConstMetric(ItemConsumptionCapacityPercent, prometheus.GaugeValue, d.ConsPercent, d.ItemName)
		ch <- prometheus.MustNewConstMetric(ItemProductionCapacityPerMinute, prometheus.GaugeValue, d.MaxProd, d.ItemName)
		ch <- prometheus.MustNewConstMetric(ItemConsumptionCapacityPerMinute, prometheus.GaugeValue, d.MaxConsumed, d.ItemName)
	}
}
