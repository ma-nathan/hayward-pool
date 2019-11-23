package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	INDEX_FILTER  = 1
	INDEX_LIGHTS  = 2
	INDEX_HEATER  = 3
	INDEX_CLEANER = 4

	EXPECTED_INDICES = 5

	STR_FILTER_OFF  = "D"
	STR_FILTER_ON   = "E"
	STR_LIGHTS_OFF  = "C"
	STR_LIGHTS_ON   = "S"
	STR_CLEANER_OFF = "C"
	STR_CLEANER_ON  = "S"
	STR_HEATER_ON   = "T"
	STR_HEATER_OFF  = "D"
)

func standardize_whitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Figure out what data we're dealing with by matching strings
// This isn't fun but it's all we have to work with

func parse_and_update(payload string) {

	var work_str string

	re := regexp.MustCompile("\r?\n")
	payload = re.ReplaceAllString(payload, "")

	re_cleanup := regexp.MustCompile("(?m)<body>(.*)</body>")

	if len(re_cleanup.FindStringSubmatch(payload)) != 2 {

		fmt.Printf("Can't parse, skipping:\n%s\n", payload)
		return
	}

	work_str = standardize_whitespace(re_cleanup.FindStringSubmatch(payload)[1])

	// fmt.Printf("Working with: %s\n", work_str)

	// For each possible string in the LCD stream, see if we can match and extract its vaule

	re = regexp.MustCompile("^Air Temp (\\d+)")
	if len(re.FindStringSubmatch(work_str)) == 2 {

		pool.AirTempF.Reading, _ = strconv.Atoi(re.FindStringSubmatch(work_str)[1])
		pool.AirTempF.Last = time.Now()
	}

	re = regexp.MustCompile("^Salt Level \\w+ (\\d+) PPM")
	if len(re.FindStringSubmatch(work_str)) == 2 {

		pool.SaltPPM.Reading, _ = strconv.Atoi(re.FindStringSubmatch(work_str)[1])
		pool.SaltPPM.Last = time.Now()
	}

	re = regexp.MustCompile("^Filter Speed \\w+ (\\w+) ")
	if len(re.FindStringSubmatch(work_str)) == 2 {

		if re.FindStringSubmatch(work_str)[1] == "Off" {

			pool.FilterSpeedRPM.Reading = 0
		} else {

			pool.FilterSpeedRPM.Reading, _ = strconv.Atoi(strings.Replace(re.FindStringSubmatch(work_str)[1], "RPM", "", -1))
		}

		pool.FilterSpeedRPM.Last = time.Now()
	}

	// Filter: ON gives us additionally:
	// Pool temp
	// Pool Chlorinator
	// The heater will turn on automatically if its temp is below pool temp (and not OFF)

	re = regexp.MustCompile("^Pool Chlorinator \\w+ (\\d+)%")
	if len(re.FindStringSubmatch(work_str)) == 2 {

		pool.ChlorinatorPct.Reading, _ = strconv.Atoi(re.FindStringSubmatch(work_str)[1])
		pool.ChlorinatorPct.Last = time.Now()
	}

	re = regexp.MustCompile("^Pool Temp (\\d+)&")
	if len(re.FindStringSubmatch(work_str)) == 2 {

		pool.PoolTempF.Reading, _ = strconv.Atoi(re.FindStringSubmatch(work_str)[1])
		pool.PoolTempF.Last = time.Now()
	}

	// AUX2:ON ("cleaner") doesn't give us any fields that Filter:ON doesn't

	// Status table:
	// filter:on aux2:on lights:off heater:on	| xxxTECT4S333333xxx
	// filter:on aux2:on lights:off	heater:off	| xxxTECD4S333333xxx
	// filter:on aux2:on lights:on				| xxxTESD4S333333xxx
	// filter:on aux2:off lights:on				| xxxTESD4C333333xxx
	// filter:off aux2:off lights:on			| xxxTDSD4C333333xxx
	// filter:off aux2:off lights:off			| xxxTDCD4C333333xxx

	// TD..4.=filter-off, TE..4.=filter-on
	// T.C.4.=lights-off, T.S.4.=lights-on
	// T...4C=cleaner-off, T...4S=cleaner-on

	re = regexp.MustCompile("xxxT(.)(.)([DT])4(.)......xxx")
	if len(re.FindStringSubmatch(work_str)) == EXPECTED_INDICES {

		switch re.FindStringSubmatch(work_str)[INDEX_FILTER] {

		case STR_FILTER_OFF:
			report_if_change(pool.FilterOn.Reading, 0, "Filter")
			pool.FilterOn.Reading = 0
			pool.FilterOn.Last = time.Now()
		case STR_FILTER_ON:
			report_if_change(pool.FilterOn.Reading, 1, "Filter")
			pool.FilterOn.Reading = 1
			pool.FilterOn.Last = time.Now()
		}

		switch re.FindStringSubmatch(work_str)[INDEX_LIGHTS] {

		case STR_LIGHTS_OFF:
			report_if_change(pool.LightOn.Reading, 0, "Lights")
			pool.LightOn.Reading = 0
			pool.LightOn.Last = time.Now()
		case STR_LIGHTS_ON:
			report_if_change(pool.LightOn.Reading, 1, "Lights")
			pool.LightOn.Reading = 1
			pool.LightOn.Last = time.Now()
		}

		switch re.FindStringSubmatch(work_str)[INDEX_CLEANER] {

		case STR_CLEANER_OFF:
			report_if_change(pool.CleanerOn.Reading, 0, "Cleaner")
			pool.CleanerOn.Reading = 0
			pool.CleanerOn.Last = time.Now()
		case STR_CLEANER_ON:
			report_if_change(pool.CleanerOn.Reading, 1, "Cleaner")
			pool.CleanerOn.Reading = 1
			pool.CleanerOn.Last = time.Now()
		}

		switch re.FindStringSubmatch(work_str)[INDEX_HEATER] {

		case STR_HEATER_OFF:
			report_if_change(pool.HeaterOn.Reading, 0, "Heater")
			pool.HeaterOn.Reading = 0
			pool.HeaterOn.Last = time.Now()
		case STR_HEATER_ON:
			report_if_change(pool.HeaterOn.Reading, 1, "Heater")
			pool.HeaterOn.Reading = 1
			pool.HeaterOn.Last = time.Now()
		}
	}
}

func report_if_change(old, new int, var_name string) {

	if old != new {

		t := time.Now()

		switch new {
		case 0:
			fmt.Printf("%s OFF at %s\n", var_name, t.Format("2006-01-02 15:04:05"))
		case 1:
			fmt.Printf("%s ON at %s\n", var_name, t.Format("2006-01-02 15:04:05"))
		}
	}
}
