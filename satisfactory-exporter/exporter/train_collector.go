package exporter

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var MaxTrainPowerConsumption = 110.0

type TrainCollector struct {
	frmTarget string
	logger    log.Logger
}

type TrainCar struct {
	Name           string  `json:"Name"`
	TotalMass      float64 `json:"TotalMass"`
	PayloadMass    float64 `json:"PayloadMass"`
	MaxPayloadMass float64 `json:"MaxPayloadMass"`
}

type TrainDetails struct {
	TrainName       string     `json:"Name"`
	PowerConsumed   float64    `json:"PowerConsumed"`
	TrainStation    string     `json:"TrainStation"`
	Derailed        bool       `json:"Derailed"`
	Status          string     `json:"Status"` //"Self-Driving",
	TrainConsist    []TrainCar `json:"TrainConsist"`
	ForwardSpeed    float64    `json:"ForwardSpeed"`
	ThrottlePercent float64    `json:"ThrottlePercent"`
	TotalMass       float64    `json:"TotalMass"`
	PayloadMass     float64    `json:"PayloadMass"`
	MaxPayloadMass  float64    `json:"MaxPayloadMass"`
}

func NewTrainCollector(frmApiAddress string, logger log.Logger) *TrainCollector {
	return &TrainCollector{
		frmTarget: frmApiAddress + "/getTrains",
		logger:    logger,
	}
}

func (c TrainCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *TrainCollector) Collect(ch chan<- prometheus.Metric) {
	details := []TrainDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading trains statistics from Ficsit Metrics", "err", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, d := range details {
		locomotives := 0.0

		for _, car := range d.TrainConsist {
			if car.Name == "Electric Locomotive" {
				locomotives = locomotives + 1
			}
		}

		// for now, the total power consumed is a multiple of the reported power consumed by the number of locomotives
		trainPowerConsumed := d.PowerConsumed * locomotives

		ch <- prometheus.MustNewConstMetric(TrainPower, prometheus.GaugeValue, trainPowerConsumed, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainTotalMass, prometheus.GaugeValue, d.TotalMass, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainPayloadMass, prometheus.GaugeValue, d.PayloadMass, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainMaxPayloadMass, prometheus.GaugeValue, d.MaxPayloadMass, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainDerailed, prometheus.GaugeValue, parseBool(d.Derailed), d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainForwardSpeed, prometheus.GaugeValue, d.ForwardSpeed, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainThrottlePercent, prometheus.GaugeValue, d.ThrottlePercent, d.TrainName)
		ch <- prometheus.MustNewConstMetric(TrainLocomotives, prometheus.GaugeValue, locomotives, d.TrainName)

		switch d.Status {
		case "Parked":
			ch <- prometheus.MustNewConstMetric(TrainDrivingStatus, prometheus.GaugeValue, 0, d.TrainName)
		case "Manual Driving":
			ch <- prometheus.MustNewConstMetric(TrainDrivingStatus, prometheus.GaugeValue, 1, d.TrainName)
		case "Self-Driving":
			ch <- prometheus.MustNewConstMetric(TrainDrivingStatus, prometheus.GaugeValue, 2, d.TrainName)
		default:
			ch <- prometheus.MustNewConstMetric(TrainDrivingStatus, prometheus.GaugeValue, -1, d.TrainName)
			c.logger.Log("msg", "Unknown train status", "status", d.Status)
		}

	}
	for circuitId, powerConsumed := range powerInfo {
		ch <- prometheus.MustNewConstMetric(TrainCircuitPower, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		ch <- prometheus.MustNewConstMetric(TrainCircuitPowerMax, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
}
