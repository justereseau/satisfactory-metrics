package exporter

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
)

type PowerInfo struct {
	CircuitId     float64 `json:"ID"`
	PowerConsumed float64 `json:"PowerConsumed"`
}

type PowerCollector struct {
	frmTarget string
	logger    log.Logger
}

type PowerDetails struct {
	CircuitId           float64 `json:"CircuitID"`
	PowerConsumed       float64 `json:"PowerConsumed"`
	PowerCapacity       float64 `json:"PowerCapacity"`
	PowerMaxConsumed    float64 `json:"PowerMaxConsumed"`
	BatteryDifferential float64 `json:"BatteryDifferential"`
	BatteryPercent      float64 `json:"BatteryPercent"`
	BatteryCapacity     float64 `json:"BatteryCapacity"`
	BatteryTimeEmpty    string  `json:"BatteryTimeEmpty"`
	BatteryTimeFull     string  `json:"BatteryTimeFull"`
	FuseTriggered       bool    `json:"FuseTriggered"`
}

func NewPowerCollector(frmApiAddress string, logger log.Logger) *PowerCollector {
	return &PowerCollector{
		frmTarget: frmApiAddress + "/getPower",
		logger:    logger,
	}
}
func (c PowerCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *PowerCollector) Collect(ch chan<- prometheus.Metric) {
	details := []PowerDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading power statistics from Ficsit Metrics", "err", err)
		return
	}

	for _, d := range details {
		circuitId := strconv.FormatFloat(d.CircuitId, 'f', -1, 64)
		ch <- prometheus.MustNewConstMetric(PowerConsumed, prometheus.GaugeValue, d.PowerConsumed, circuitId)
		ch <- prometheus.MustNewConstMetric(PowerCapacity, prometheus.GaugeValue, d.PowerCapacity, circuitId)
		ch <- prometheus.MustNewConstMetric(PowerMaxConsumed, prometheus.GaugeValue, d.PowerMaxConsumed, circuitId)
		ch <- prometheus.MustNewConstMetric(BatteryDifferential, prometheus.GaugeValue, d.BatteryDifferential, circuitId)
		ch <- prometheus.MustNewConstMetric(BatteryPercent, prometheus.GaugeValue, d.BatteryPercent, circuitId)
		ch <- prometheus.MustNewConstMetric(BatteryCapacity, prometheus.GaugeValue, d.BatteryCapacity, circuitId)
		batterySecondsRemaining := parseTimeSeconds(d.BatteryTimeEmpty)
		if batterySecondsRemaining != nil {
			ch <- prometheus.MustNewConstMetric(BatterySecondsEmpty, prometheus.GaugeValue, *batterySecondsRemaining, circuitId)
		}
		batterySecondsFull := parseTimeSeconds(d.BatteryTimeFull)
		if batterySecondsFull != nil {
			ch <- prometheus.MustNewConstMetric(BatterySecondsFull, prometheus.GaugeValue, *batterySecondsFull, circuitId)
		}
		ch <- prometheus.MustNewConstMetric(FuseTriggered, prometheus.GaugeValue, parseBool(d.FuseTriggered), circuitId)
	}
}
