package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	BASE_URL          = "http://bbq.iot.fumanchu.com:8080"
	POOL_DEVICE_TOKEN = "diCFzJhC2pXp1M0rqWnv"
	TB_HTTP_TIMEOUT   = 10
)

func add_json_element(in, in_name string, in_val int) (out string) {

	// This lead to weird graph tails instead of the gaps I was expecting

/*
	if in_val == NOT_RECORDED {
		out = in
		return
	}
*/

	out = in + "\"" + in_name + "\":" + strconv.Itoa(in_val) + ","
	return
}

func http_call_thingsboard(json_str string) {

	var resp *http.Response
	var req *http.Request

	url := BASE_URL + "/api/v1/" + POOL_DEVICE_TOKEN + "/telemetry"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json_str)))

	if err != nil {

		fmt.Printf("http_call_thingsboard: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{} // {Timeout: TB_HTTP_TIMEOUT}
	resp, err = client.Do(req)

	if err != nil {

		fmt.Printf("http_call_thingsboard: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.Status != "200 " {

		fmt.Printf("Status: \"%s\"\n", resp.Status)
		fmt.Println("Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Body:", string(body))
		fmt.Printf("We sent: %s\n", json_str)
	}
}

func deliver_stats_to_thingsboard() {

	var json_str string

	json_str = "{"
	json_str = add_json_element(json_str, "Air Temp", pool.AirTempF.Reading)
	json_str = add_json_element(json_str, "Pool Temp", pool.PoolTempF.Reading)
	json_str = add_json_element(json_str, "Filter Speed", pool.FilterSpeedRPM.Reading)
	json_str = add_json_element(json_str, "Salt", pool.SaltPPM.Reading)
	json_str = add_json_element(json_str, "Filter", pool.FilterOn.Reading)
	json_str = add_json_element(json_str, "Cleaner", pool.CleanerOn.Reading)
	json_str = add_json_element(json_str, "Lights", pool.LightOn.Reading)
	json_str = add_json_element(json_str, "Heater", pool.HeaterOn.Reading)
	json_str = add_json_element(json_str, "Chlorinator", pool.ChlorinatorPct.Reading)
	json_str = json_str + "}"

	json_str = strings.Replace(json_str, ",}", "}", -1) // Final fixup

	fmt.Printf("Payload: %s\n", json_str)

	http_call_thingsboard(json_str)
	//	fmt.Printf("deliver_stats_to_thingsboard: ok\n")
}
