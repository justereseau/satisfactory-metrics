package exporter

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type PlayerCollector struct {
	frmTarget string
	logger    log.Logger
}

type TagColor struct {
	R float64 `json:"R"`
	G float64 `json:"G"`
	B float64 `json:"B"`
	A float64 `json:"A"`
}

type PlayerDetails struct {
	ID         float64  `json:"ID"`
	PlayerName string   `json:"PlayerName"`
	PlayerHP   float64  `json:"PlayerHP"`
	Dead       bool     `json:"Dead"`
	PingTime   float64  `json:"PingTime"`
	Location   Location `json:"Location"`
	TagColor   TagColor `json:"TagColor"`
}

func NewPlayerCollector(frmApiAddress string, logger log.Logger) *PlayerCollector {
	return &PlayerCollector{
		frmTarget: frmApiAddress + "/getPlayer",
		logger:    logger,
	}
}

func (c PlayerCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *PlayerCollector) Collect(ch chan<- prometheus.Metric) {
	details := []PlayerDetails{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading trains statistics from Ficsit Metrics", "err", err)
		return
	}

	for _, d := range details {
		ch <- prometheus.MustNewConstMetric(PlayerPosition, prometheus.GaugeValue, d.Location.X*0.000000117118912+0.03804908, d.PlayerName, fmt.Sprintf("%f", d.ID), "X")
		ch <- prometheus.MustNewConstMetric(PlayerPosition, prometheus.GaugeValue, d.Location.Y*0.000000117118912-0.0439383731, d.PlayerName, fmt.Sprintf("%f", d.ID), "Y")
		ch <- prometheus.MustNewConstMetric(PlayerPosition, prometheus.GaugeValue, d.Location.Z, d.PlayerName, fmt.Sprintf("%f", d.ID), "Z")
		ch <- prometheus.MustNewConstMetric(PlayerRotation, prometheus.GaugeValue, float64(d.Location.Rotation), d.PlayerName, fmt.Sprintf("%f", d.ID))
		ch <- prometheus.MustNewConstMetric(PlayerHealth, prometheus.GaugeValue, d.PlayerHP, d.PlayerName, fmt.Sprintf("%f", d.ID))
		ch <- prometheus.MustNewConstMetric(PlayerDead, prometheus.GaugeValue, parseBool(d.Dead), d.PlayerName, fmt.Sprintf("%f", d.ID))
		ch <- prometheus.MustNewConstMetric(PlayerPing, prometheus.GaugeValue, d.PingTime, d.PlayerName, fmt.Sprintf("%f", d.ID))
		ch <- prometheus.MustNewConstMetric(PlayerTagColor, prometheus.GaugeValue, d.TagColor.R, d.PlayerName, fmt.Sprintf("%f", d.ID), "R")
		ch <- prometheus.MustNewConstMetric(PlayerTagColor, prometheus.GaugeValue, d.TagColor.G, d.PlayerName, fmt.Sprintf("%f", d.ID), "G")
		ch <- prometheus.MustNewConstMetric(PlayerTagColor, prometheus.GaugeValue, d.TagColor.B, d.PlayerName, fmt.Sprintf("%f", d.ID), "B")
		ch <- prometheus.MustNewConstMetric(PlayerTagColor, prometheus.GaugeValue, d.TagColor.A, d.PlayerName, fmt.Sprintf("%f", d.ID), "A")
	}
}
