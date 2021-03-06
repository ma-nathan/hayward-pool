package main

import (
	//	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

// CREATE USER admin WITH PASSWORD '$the_usual' WITH ALL PRIVILEGES
// create database BLAH

func influxDBClient(config Config) client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.DatabaseURL,
		Username: config.DatabaseUser,
		Password: config.DatabasePassword,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}

func influx_push_metrics(c client.Client, config Config) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.DatabaseDatabase,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	eventTime := time.Now()

	/*
		Using "Line Protocol", eg: cpu,host=server02,region=uswest value=3 1434055562000010000
		http://goinbigdata.com/working-with-influxdb-in-go/

		key: pool
		tags: none
		fields: pool_temp=blah, etc.
		timestamp in seconds
	*/

	key := "pool"
	tags := map[string]string{}
	fields := map[string]interface{}{
		"air_temp": pool.AirTempF.Reading,
	}

	point, err := client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	if pool.AirTempF.Reading != NOT_RECORDED {
		bp.AddPoint(point)
	}

	fields = map[string]interface{}{
		"pool_temp": pool.PoolTempF.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	if pool.PoolTempF.Reading != NOT_RECORDED {
		bp.AddPoint(point)
	}

	fields = map[string]interface{}{
		"filter_speed": pool.FilterSpeedRPM.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	if pool.FilterSpeedRPM.Reading != NOT_RECORDED {
		bp.AddPoint(point)
	}

	fields = map[string]interface{}{
		"salt_ppm": pool.SaltPPM.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	if pool.SaltPPM.Reading != NOT_RECORDED {
		bp.AddPoint(point)
	}

	fields = map[string]interface{}{
		"filter_on": pool.FilterOn.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	bp.AddPoint(point)

	fields = map[string]interface{}{
		"cleaner_on": pool.CleanerOn.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	bp.AddPoint(point)

	fields = map[string]interface{}{
		"lights_on": pool.LightOn.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	bp.AddPoint(point)

	fields = map[string]interface{}{
		"chlorinator_percent": pool.ChlorinatorPct.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	if pool.ChlorinatorPct.Reading != NOT_RECORDED {
		bp.AddPoint(point)
	}

	fields = map[string]interface{}{
		"heater_on": pool.HeaterOn.Reading,
	}

	point, err = client.NewPoint(key, tags, fields, eventTime)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	bp.AddPoint(point)

	err = c.Write(bp)
	if err != nil {
		log.Fatal(err)
	}

}

func deliver_stats_to_influxdb(c client.Client, config Config) {

	influx_push_metrics(c, config)
}
