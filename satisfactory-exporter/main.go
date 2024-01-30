package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/go-kit/log/level"
	"github.com/justereseau/satisfactory-metrics/satisfactory-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
)

// Define constants
const (
	exporter_name         = "ficsit_remote_monitoring"
	exporter_display_name = "Ficsit Remote Monitoring Exporter"
)

// Define parameters
var (
	listenAddress = flag.String("web.listen-address", "127.0.0.1:9100", "Address to listen on for web interface and telemetry.")
	logLevel      = flag.String("log.level", "info", "Only log messages with the given severity or above. One of: [debug, info, warn, error, none]")
	frmApiAddress = flag.String("frm.listen-address", "http://localhost:8080", "Address of Ficsit Remote Monitoring webserver")
)

func main() {
	// Get parameters
	flag.Parse()

	allowedLogLevel := &promlog.AllowedLevel{}
	allowedLogLevel.Set(*logLevel)
	logger := promlog.New(&promlog.Config{Level: allowedLogLevel})

	level.Info(logger).Log("msg", "Starting "+exporter_name+"_exporter.")

	prometheus.MustRegister(version.NewCollector(exporter_name + "_exporter"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>` + exporter_display_name + `</title></head>
			<body>
			<h1>` + exporter_display_name + `</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		level.Debug(logger).Log("msg", "Starting scrape")

		registry := prometheus.NewRegistry()
		registry.MustRegister(exporter.NewProductionCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewPowerCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewFactoryBuildingCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewVehicleCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewDroneStationCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewVehicleStationCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewTrainCollector(*frmApiAddress, logger))
		registry.MustRegister(exporter.NewTrainStationCollector(*frmApiAddress, logger))

		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		level.Debug(logger).Log("msg", "Scrape done.", "duration", time.Since(start).Seconds())
	})

	level.Info(logger).Log("msg", "Starting to listen.", "address", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to start http server.", "err", err)
	}
}
