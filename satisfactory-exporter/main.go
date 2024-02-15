package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
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

	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ok`))
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		level.Debug(logger).Log("msg", "Starting scrape")

		registry := prometheus.NewRegistry()

		// Get enabled collectors from request
		enabledCollectors := r.URL.Query().Get("collect")
		fmt.Println("Enabled collectors: ", enabledCollectors)
		if enabledCollectors == "all" || enabledCollectors == "" {
			enabledCollectors = "production,power,factory_building,vehicle,drone_station,vehicle_station,train,train_station,player"
		}
		for _, collector := range strings.Split(enabledCollectors, ",") {
			fmt.Println("Registering collector: ", collector)
			switch collector {
			case "production":
				registry.MustRegister(exporter.NewProductionCollector(*frmApiAddress, logger))
			case "power":
				registry.MustRegister(exporter.NewPowerCollector(*frmApiAddress, logger))
			case "factory_building":
				registry.MustRegister(exporter.NewFactoryBuildingCollector(*frmApiAddress, logger))
			case "vehicle":
				registry.MustRegister(exporter.NewVehicleCollector(*frmApiAddress, logger))
			case "drone_station":
				registry.MustRegister(exporter.NewDroneStationCollector(*frmApiAddress, logger))
			case "vehicle_station":
				registry.MustRegister(exporter.NewVehicleStationCollector(*frmApiAddress, logger))
			case "train":
				registry.MustRegister(exporter.NewTrainCollector(*frmApiAddress, logger))
			case "train_station":
				registry.MustRegister(exporter.NewTrainStationCollector(*frmApiAddress, logger))
			case "player":
				registry.MustRegister(exporter.NewPlayerCollector(*frmApiAddress, logger))
			default:
				level.Warn(logger).Log("msg", "Unknown collector", "collector", collector)
			}
		}

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
