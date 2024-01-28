package exporter

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	StationPower       = 0.1 // should be 50, but currently bugged.
	CargoPlatformPower = 50.0
)

type TrainStationCollector struct {
	frmTarget string
	logger    log.Logger
}

type CargoPlatform struct {
	LoadingDock   string  `json:"LoadingDock"`
	TransferRate  float64 `json:"TransferRate"`
	LoadingStatus string  `json:"LoadingStatus"` // Idle, Loading, Unloading
	LoadingMode   string  `json:"LoadingMode"`
}

type TrainStationDetails struct {
	Name           string          `json:"Name"`
	Location       Location        `json:"location"`
	CargoPlatforms []CargoPlatform `json:"CargoPlatforms"`
	PowerInfo      PowerInfo       `json:"PowerInfo"`
}

func NewTrainStationCollector(frmApiAddress string, logger log.Logger) *TrainStationCollector {
	return &TrainStationCollector{
		frmTarget: frmApiAddress + "/getTrainStation",
		logger:    logger,
	}
}

func (c TrainStationCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *TrainStationCollector) Collect(ch chan<- prometheus.Metric) {
	details := []TrainStationDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading train station statistics from Ficsit Metrics", "err", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		val, ok := powerInfo[d.PowerInfo.CircuitId]
		maxval, maxok := maxPowerInfo[d.PowerInfo.CircuitId]

		// some additional calculations: for now, power listed here is only for the station.
		// add each of the cargo platforms' power info: 0.1MW if Idle, 50MW otherwise
		totalPowerConsumed := d.PowerInfo.PowerConsumed
		maxTotalPowerConsumed := StationPower
		for _, p := range d.CargoPlatforms {
			maxTotalPowerConsumed = maxTotalPowerConsumed + CargoPlatformPower
			if p.LoadingStatus == "Idle" {
				totalPowerConsumed = totalPowerConsumed + 0.1
			} else {
				totalPowerConsumed = totalPowerConsumed + CargoPlatformPower
			}
		}

		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + totalPowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = totalPowerConsumed
		}

		if maxok {
			maxPowerInfo[d.PowerInfo.CircuitId] = maxval + maxTotalPowerConsumed
		} else {
			maxPowerInfo[d.PowerInfo.CircuitId] = maxTotalPowerConsumed
		}
	}

	for circuitId, powerConsumed := range powerInfo {
		ch <- prometheus.MustNewConstMetric(TrainStationPower, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}

	for circuitId, powerConsumed := range maxPowerInfo {
		ch <- prometheus.MustNewConstMetric(TrainStationPowerMax, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))

	}
}
