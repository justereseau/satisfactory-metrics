package exporter

import (
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type VehicleCollector struct {
	frmTarget string
	logger    log.Logger
}

type VehicleDetails struct {
	Id           string   `json:"ID"`
	VehicleType  string   `json:"Name"`
	Location     Location `json:"location"`
	ForwardSpeed float64  `json:"ForwardSpeed"`
	AutoPilot    bool     `json:"AutoPilot"`
	Fuel         []Fuel   `json:"Fuel"`
	PathName     string   `json:"PathName"`
	DepartTime   time.Time
	Departed     bool
}

type Fuel struct {
	Name   string  `json:"Name"`
	Amount float64 `json:"Amount"`
}

func NewVehicleCollector(frmApiAddress string, logger log.Logger) *VehicleCollector {
	return &VehicleCollector{
		frmTarget: frmApiAddress + "/getVehicles",
		logger:    logger,
	}
}

func (c VehicleCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *VehicleCollector) Collect(ch chan<- prometheus.Metric) {
	details := []VehicleDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading vehicle statistics from Ficsit Metrics", "err", err)
		return
	}

	for _, d := range details {
		for n, f := range d.Fuel {
			ch <- prometheus.MustNewConstMetric(VehicleFuel, prometheus.GaugeValue, f.Amount, d.Id, d.VehicleType, f.Name, strconv.Itoa(n))
		}
	}
}
