package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PlayerPosition = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_current_position",
		Help: "The current position of the player, per axis",
	}, []string{
		"player_name",
		"player_id",
		"axis_name",
	})
	PlayerRotation = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_current_rotation",
		Help: "The current rotation of the player, in degrees",
	}, []string{
		"player_name",
		"player_id",
	})
	PlayerHealth = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_current_health",
		Help: "The current health of the player",
	}, []string{
		"player_name",
		"player_id",
	})
	PlayerDead = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_is_dead",
		Help: "Is the player dead",
	}, []string{
		"player_name",
		"player_id",
	})
	PlayerPing = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_current_ping",
		Help: "The current ping of the player",
	}, []string{
		"player_name",
		"player_id",
	})
	PlayerTagColor = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "player_tag_color",
		Help: "The current tag color of the player",
	}, []string{
		"player_name",
		"player_id",
		"component",
	})

	ItemProductionCapacityPerMinute = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_production_capacity_per_min",
		Help: "The factory's capacity for the production of an item, per minute",
	}, []string{
		"item_name",
	})

	ItemProductionCapacityPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_production_capacity_pc",
		Help: "The percentage of an item's production capacity being used",
	}, []string{
		"item_name",
	})

	ItemConsumptionCapacityPerMinute = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_consumption_capacity_per_min",
		Help: "The factory's capacity for the consumption of an item, per minute",
	}, []string{
		"item_name",
	})

	ItemConsumptionCapacityPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "item_consumption_capacity_pc",
		Help: "The percentage of an item's consumption capacity being used",
	}, []string{
		"item_name",
	})

	ItemsProducedPerMin = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "items_produced_per_min",
		Help: "The number of an item being produced, per minute",
	}, []string{
		"item_name",
	})

	ItemsConsumedPerMin = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "items_consumed_per_min",
		Help: "The number of an item being consumed, per minute",
	}, []string{
		"item_name",
	})

	PowerConsumed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_consumed",
		Help: "Power consumed on selected power circuit",
	}, []string{
		"circuit_id",
	})

	PowerCapacity = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_capacity",
		Help: "Power capacity on selected power circuit",
	}, []string{
		"circuit_id",
	})

	PowerMaxConsumed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "power_max_consumed",
		Help: "Maximum Power that can be consumed on selected power circuit",
	}, []string{
		"circuit_id",
	})

	BatteryDifferential = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_differential",
		Help: "Amount of power in excess/deficit going into or out of the battery bank(s). Positive = Charges batteries, Negative = Drains batteries",
	}, []string{
		"circuit_id",
	})

	BatteryPercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_percent",
		Help: "Percentage of battery bank(s) charge",
	}, []string{
		"circuit_id",
	})

	BatteryCapacity = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_capacity",
		Help: "Total capacity of battery bank(s)",
	}, []string{
		"circuit_id",
	})

	BatterySecondsEmpty = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_seconds_empty",
		Help: "Seconds until Batteries are empty",
	}, []string{
		"circuit_id",
	})

	BatterySecondsFull = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_seconds_full",
		Help: "Seconds until Batteries are full",
	}, []string{
		"circuit_id",
	})

	FuseTriggered = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "fuse_triggered",
		Help: "Has the fuse been triggered",
	}, []string{
		"circuit_id",
	})

	VehicleFuel = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_fuel",
		Help: "Amount of fuel remaining",
	}, []string{
		"id",
		"vehicle_type",
		"fuel_type",
		"fuel_index",
	})

	DronePortBatteryRate = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_battery_rate",
		Help: "Rate of batteries used",
	}, []string{
		"id",
		"home_station",
		"paired_station",
	})

	DronePortRndTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_round_trip_seconds",
		Help: "Recorded drone round trip time in seconds",
	}, []string{
		"id",
		"home_station",
		"paired_station",
	})

	DronePortPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "drone_port_power",
		Help: "Drone port power in MW",
	}, []string{
		"circuit_id",
	})

	VehicleStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_station_power",
		Help: "Vehicle station power use in MW",
	}, []string{
		"circuit_id",
	})
	VehicleStationPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "vehicle_station_power_max",
		Help: "Vehicle station max power use in MW",
	}, []string{
		"circuit_id",
	})

	TrainRoundTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_round_trip_seconds",
		Help: "Recorded train round trip time in seconds",
	}, []string{
		"name",
	})
	TrainSegmentTrip = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_segment_trip_seconds",
		Help: "Recorded train trip between two stations",
	}, []string{
		"name",
		"from",
		"to",
	})
	TrainDerailed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_derailed",
		Help: "Is train derailed",
	}, []string{
		"name",
	})
	TrainPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_power_consumed",
		Help: "How much power train is consuming",
	}, []string{
		"name",
	})
	TrainDrivingStatus = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_driving_status",
		Help: "The current autopilot status of the train. 0 = Manual, 1 = Autopilot",
	}, []string{
		"name",
	})
	TrainForwardSpeed = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_forward_speed",
		Help: "The current forward speed of the train",
	}, []string{
		"name",
	})
	TrainThrottlePercent = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_throttle_percent",
		Help: "The current throttle percentage of the train",
	}, []string{
		"name",
	})
	TrainLocomotives = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_locomotives",
		Help: "The number of locomotives on the train",
	}, []string{
		"name",
	})
	TrainCircuitPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_power_circuit_consumed",
		Help: "How much power all trains are consuming in a circuit",
	}, []string{
		"circuit_id",
	})
	TrainCircuitPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_power_circuit_consumed_max",
		Help: "Maximum power all trains can consume on a circuit",
	}, []string{
		"circuit_id",
	})
	TrainTotalMass = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_total_mass",
		Help: "Total mass of the train",
	}, []string{
		"name",
	})
	TrainPayloadMass = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_payload_mass",
		Help: "Current payload mass of the train",
	}, []string{
		"name",
	})
	TrainMaxPayloadMass = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_max_payload_mass",
		Help: "Max payload mass of the train",
	}, []string{
		"name",
	})

	TrainStationPower = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_station_power",
		Help: "Train station power consumed in MW",
	}, []string{
		"circuit_id",
	})
	TrainStationPowerMax = RegisterNewGaugeVec(prometheus.GaugeOpts{
		Name: "train_station_power_max",
		Help: "Train station power max consumed in MW",
	}, []string{
		"circuit_id",
	})
)
