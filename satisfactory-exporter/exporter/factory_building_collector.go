package exporter

import (
	"math"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type FactoryBuildingCollector struct {
	frmTarget string
	logger    log.Logger
}

func NewFactoryBuildingCollector(frmApiAddress string, logger log.Logger) *FactoryBuildingCollector {
	return &FactoryBuildingCollector{
		frmTarget: frmApiAddress + "/getFactory",
		logger:    logger,
	}
}

func (c FactoryBuildingCollector) Describe(ch chan<- *prometheus.Desc) {}

func (c *FactoryBuildingCollector) Collect(ch chan<- prometheus.Metric) {
	details := []BuildingDetail{}
	err := retrieveData(c.frmTarget, &details)
	if err != nil {
		level.Error(c.logger).Log("msg", "Error reading building statistics from Ficsit Metrics", "err", err)
		return
	}

	powerInfo := map[float64]float64{}
	maxPowerInfo := map[float64]float64{}
	for _, building := range details {
		for _, prod := range building.Production {
			ch <- prometheus.MustNewConstMetric(
				MachineItemsProducedPerMin,
				prometheus.GaugeValue,
				prod.CurrentProd,
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
			)

			ch <- prometheus.MustNewConstMetric(
				MachineItemsProducedEffiency,
				prometheus.GaugeValue,
				prod.ProdPercent,
				prod.Name,
				building.Building,
				strconv.FormatFloat(building.Location.X, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Y, 'f', -1, 64),
				strconv.FormatFloat(building.Location.Z, 'f', -1, 64),
			)
		}

		val, ok := powerInfo[building.PowerInfo.CircuitId]
		if ok {
			powerInfo[building.PowerInfo.CircuitId] = val + building.PowerInfo.PowerConsumed
		} else {
			powerInfo[building.PowerInfo.CircuitId] = building.PowerInfo.PowerConsumed
		}
		val, ok = maxPowerInfo[building.PowerInfo.CircuitId]
		maxBuildingPower := 0.0
		switch building.Building {
		case "Smelter":
			maxBuildingPower = SmelterPower
		case "Constructor":
			maxBuildingPower = ConstructorPower
		case "Assembler":
			maxBuildingPower = AssemblerPower
		case "Manufacturer":
			maxBuildingPower = ManufacturerPower
		case "Blender":
			maxBuildingPower = BlenderPower
		case "Refinery":
			maxBuildingPower = RefineryPower
		case "Particle Accelerator":
			maxBuildingPower = ParticleAcceleratorPower
		}
		//update max power from clock speed
		// see https://satisfactory.wiki.gg/wiki/Clock_speed#Clock_speed_for_production_buildings for power info
		maxBuildingPower = maxBuildingPower * (math.Pow(building.ManuSpeed/100, 1.321928))
		if ok {
			maxPowerInfo[building.PowerInfo.CircuitId] = val + maxBuildingPower
		} else {
			maxPowerInfo[building.PowerInfo.CircuitId] = maxBuildingPower
		}
	}
	for circuitId, powerConsumed := range powerInfo {
		ch <- prometheus.MustNewConstMetric(FactoryPower, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
	for circuitId, powerConsumed := range maxPowerInfo {
		ch <- prometheus.MustNewConstMetric(FactoryPowerMax, prometheus.GaugeValue, powerConsumed, strconv.FormatFloat(circuitId, 'f', -1, 64))
	}
}
