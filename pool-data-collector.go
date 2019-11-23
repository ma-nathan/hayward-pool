package main

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ASSUME_GONE    = -1 * time.Minute
	ENDPOINT_PAUSE = time.Second * 2
	HTTP_TIMEOUT   = time.Second * 30
	DATA_UPDATE    = time.Minute * 2
	NOT_RECORDED   = 0
	version        = "0.1"
)

var pool PoolData

/*
	The heater is now wired to the controller 11/22/2019 after 8+ years.

*/

func get_lcd_payload(url string) (payload string, err error) {

	var resp *http.Response
	var req *http.Request
	var http_err error
	var data []byte

	client := &http.Client{Timeout: HTTP_TIMEOUT}
	req, http_err = http.NewRequest("GET", url, nil)

	if http_err != nil {
		return "", http_err
	}

	resp, http_err = client.Do(req)

	if http_err != nil {
		return "", http_err
	}

	defer resp.Body.Close()
	data, http_err = ioutil.ReadAll(resp.Body)

	payload = string(data)
	return
}

func watch_http_endpoint(config Config) {

	// Treat the "LCD display" like a serial endpoint over which we have no control on
	// the sending side.  Keep polling to see what it currently has to say and update
	// our tracking to match.  Availability will come and go.

	for {

		payload, err := get_lcd_payload("http://" + config.PoolHost + "/WNewSt.htm")

		if err != nil {

			fmt.Printf("Error fetching data from HTTP endpoint: %v\n", err)

		} else {

			// Send it over to get parsed

			parse_and_update(payload)
		}

		time.Sleep(ENDPOINT_PAUSE)
	}
}

func update_datastore(c client.Client, config Config) {

	// Every DATA_UPDATE interval:
	// 1. Check if our data is stale, zero it out if so
	// 2. Write what we have to datastore

	for {

		time.Sleep(DATA_UPDATE) // don't deliver first thing before we have data

		if pool.AirTempF.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.AirTempF.Reading = NOT_RECORDED
		}

		if pool.PoolTempF.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.PoolTempF.Reading = NOT_RECORDED
		}

		if pool.FilterSpeedRPM.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.FilterSpeedRPM.Reading = NOT_RECORDED
		}

		if pool.SaltPPM.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.SaltPPM.Reading = NOT_RECORDED
		}

		if pool.FilterOn.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.FilterOn.Reading = NOT_RECORDED
		}

		if pool.CleanerOn.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.CleanerOn.Reading = NOT_RECORDED
		}

		if pool.LightOn.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.LightOn.Reading = NOT_RECORDED
		}

		if pool.ChlorinatorPct.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.ChlorinatorPct.Reading = NOT_RECORDED
		}

		if pool.HeaterOn.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.HeaterOn.Reading = NOT_RECORDED
		}

		/*
			fmt.Printf("AirTempF: %d\n", pool.AirTempF.Reading)
			fmt.Printf("PoolTempF: %d\n", pool.PoolTempF.Reading)
			fmt.Printf("FilterSpeedRPM: %d\n", pool.FilterSpeedRPM.Reading)
			fmt.Printf("SaltPPM: %d\n", pool.SaltPPM.Reading)
			fmt.Printf("ChlorinatorPct: %d\n", pool.ChlorinatorPct.Reading)
			fmt.Printf("FilterOn: %d\n", pool.FilterOn.Reading)
			fmt.Printf("CleanerOn: %d\n", pool.CleanerOn.Reading)
			fmt.Printf("LightOn: %d\n", pool.LightOn.Reading)
		*/

		// Now deliver this data to the backend
		// We can support thingsboard, kairosdb, and influxdb

		//deliver_stats_to_kairos()
		deliver_stats_to_influxdb(c, config)

	}

}
