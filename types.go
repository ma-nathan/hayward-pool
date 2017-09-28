package main

import (
	"time"
)

type Measurement struct {
	Reading int
	Last    time.Time
}

type PoolData struct {
	AirTempF       Measurement
	PoolTempF      Measurement
	FilterSpeedRPM Measurement
	SaltPPM        Measurement
	FilterOn       Measurement
	CleanerOn      Measurement
	LightOn        Measurement
	HeaterOn       Measurement
	ChlorinatorPct Measurement
}
