package exporter

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var VehicleStationPowerConsumption = 20.0

type VehicleStationCollector struct {
	frmTarget string
	logger    log.Logger
}

type VehicleStationDetails struct {
	Name      string    `json:"Name"`
	Location  Location  `json:"location"`
	PowerInfo PowerInfo `json:"PowerInfo"`
}

func NewVehicleStationCollector(frmApiAddress string, logger log.Logger) *VehicleStationCollector {
	return &VehicleStationCollector{
		frmTarget: frmApiAddress + "/getTruckStation",
		logger:    logger,
	}
}

func (c VehicleStationCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *VehicleStationCollector) Collect(ch chan<- prometheus.Metric) {
	details := []VehicleStationDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading vehicle station statistics from Ficsit Metrics", "err", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitId]
		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + d.PowerInfo.PowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = d.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[d.PowerInfo.CircuitId]
		if ok {
			maxPowerInfo[d.PowerInfo.CircuitId] = val + VehicleStationPowerConsumption
		} else {
			maxPowerInfo[d.PowerInfo.CircuitId] = VehicleStationPowerConsumption
		}
	}

	for circuitId, powerConsumed := range powerInfo {
		ch <- prometheus.MustNewConstMetric(VehicleStationPower, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}

	for circuitId, powerConsumed := range maxPowerInfo {
		ch <- prometheus.MustNewConstMetric(VehicleStationPowerMax, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
}
