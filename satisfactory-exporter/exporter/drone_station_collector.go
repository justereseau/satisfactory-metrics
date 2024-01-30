package exporter

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type DroneStationCollector struct {
	frmTarget string
	logger    log.Logger
}

type DroneStationDetails struct {
	Id                     string    `json:"ID"`
	HomeStation            string    `json:"Name"`
	PairedStation          string    `json:"PairedStation"`
	DroneStatus            string    `json:"DroneStatus"`
	AvgIncRate             float64   `json:"AvgIncRate"`
	AvgIncStack            float64   `json:"AvgIncStack"`
	AvgOutRate             float64   `json:"AvgOutRate"`
	AvgOutStack            float64   `json:"AvgOutStack"`
	AvgRndTrip             string    `json:"AvgRndTrip"`
	AvgTotalIncRate        float64   `json:"AvgTotalIncRate"`
	AvgTotalIncStack       float64   `json:"AvgTotalIncStack"`
	AvgTotalOutRate        float64   `json:"AvgTotalOutRate"`
	AvgTotalOutStack       float64   `json:"AvgTotalOutStack"`
	AvgTripIncAmt          float64   `json:"AvgTripIncAmt"`
	EstRndTrip             string    `json:"EstRndTrip"`
	EstTotalTransRate      float64   `json:"EstTotalTransRate"`
	EstTransRate           float64   `json:"EstTransRate"`
	EstLatestTotalIncStack float64   `json:"EstLatestTotalIncStack"`
	EstLatestTotalOutStack float64   `json:"EstLatestTotalOutStack"`
	LatestIncStack         float64   `json:"LatestIncStack"`
	LatestOutStack         float64   `json:"LatestOutStack"`
	LatestRndTrip          string    `json:"LatestRndTrip"`
	LatestTripIncAmt       float64   `json:"LatestTripIncAmt"`
	LatestTripOutAmt       float64   `json:"LatestTripOutAmt"`
	MedianRndTrip          string    `json:"MedianRndTrip"`
	MedianTripIncAmt       float64   `json:"MedianTripIncAmt"`
	MedianTripOutAmt       float64   `json:"MedianTripOutAmt"`
	EstBatteryRate         float64   `json:"EstBatteryRate"`
	PowerInfo              PowerInfo `json:"PowerInfo"`
}

func NewDroneStationCollector(frmApiAddress string, logger log.Logger) *DroneStationCollector {
	return &DroneStationCollector{
		frmTarget: frmApiAddress + "/getDroneStation",
		logger:    logger,
	}
}

func (c DroneStationCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *DroneStationCollector) Collect(ch chan<- prometheus.Metric) {
	details := []DroneStationDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading drone statistics from Ficsit Metrics", "err", err)
		return
	}

	powerInfo := map[float64]float64{}
	for _, d := range details {
		id := d.Id
		home := d.HomeStation
		paired := d.PairedStation

		ch <- prometheus.MustNewConstMetric(DronePortBatteryRate, prometheus.GaugeValue, d.EstBatteryRate, id, home, paired)

		roundTrip := parseTimeSeconds(d.LatestRndTrip)
		if roundTrip != nil {
			ch <- prometheus.MustNewConstMetric(DronePortRndTrip, prometheus.GaugeValue, *roundTrip, id, home, paired)
		}

		val, ok := powerInfo[d.PowerInfo.CircuitId]
		if ok {
			powerInfo[d.PowerInfo.CircuitId] = val + d.PowerInfo.PowerConsumed
		} else {
			powerInfo[d.PowerInfo.CircuitId] = d.PowerInfo.PowerConsumed
		}
	}

	for circuitId, powerConsumed := range powerInfo {
		ch <- prometheus.MustNewConstMetric(DronePortPower, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
}
