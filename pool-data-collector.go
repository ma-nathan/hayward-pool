package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DEVICE_HOST      = "pool.fumanchu.com"
	ASSUME_GONE      = -1 * time.Minute
	ENDPOINT_PAUSE   = time.Second * 2
	HTTP_TIMEOUT     = time.Second * 30
	DATA_UPDATE      = time.Minute * 1
	POOL_TEMP_TARGET = 88
	NOT_RECORDED     = 0
	// NOT_RECORDED     = -1
	version = "0.1"
)

var pool PoolData

/*
	FIXME: need to make some adjustments to inference of state of heater

	"Cleaner" AKA "Pool pump" is AUX2 (just FYI in case you're using the UI)

	We infer that the heater is ON when:
	- Pool temp is less than POOL_TEMP_TARGET
	- "Cleaner" / "Pool pump" is ON

	This is wrong because:
	- Heater runs when Filter is ON, not tied to "Cleaner" / "Pool pump"
	- POOL_TEMP_TARGET is hardcoded and is completely wrong for winter

	Looking at the code, that's not actually what we do, but...it's how
	the graphs react.  So, investigate what's going on.

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

func watch_http_endpoint() {

	// Treat the "LCD display" like a serial endpoint over which we have no control on
	// the sending side.  Keep polling to see what it currently has to say and update
	// our tracking to match.  Availability will come and go.

	for {

		payload, err := get_lcd_payload("http://" + DEVICE_HOST + "/WNewSt.htm")

		if err != nil {

			fmt.Printf("Error fetching data from HTTP endpoint: %v\n", err)

		} else {

			// Send it over to get parsed

			parse_and_update(payload)
		}

		time.Sleep(ENDPOINT_PAUSE)
	}
}

func update_datastore( target_temp int ) {

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

		// Special case: infer state of heater

		if pool.FilterOn.Reading == 1 && pool.PoolTempF.Reading < target_temp {

			pool.HeaterOn.Reading = 1
			pool.HeaterOn.Last = time.Now()
		} else {

			pool.HeaterOn.Reading = 0
		}

		if pool.FilterOn.Last.Before(time.Now().Add(ASSUME_GONE)) {
			pool.HeaterOn.Reading = NOT_RECORDED
		}

		// Now deliver this data to the backend - thingsboard?  cassandra?
		// We can support thingsboard and kairosdb

		deliver_stats_to_kairos()

	}

}

func main() {

	fmt.Println("pool-data-collector polls a Hayward Aqua Connect Local network device.")
	fmt.Println("Data is uploaded to a thingsboard/kairosdb instance for graphing.")

	target_temp := handle_command_line_args()

	go update_datastore( target_temp )
	watch_http_endpoint()
}
