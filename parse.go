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
	INDEX_CLEANER = 3

	STR_FILTER_OFF  = "D"
	STR_FILTER_ON   = "E"
	STR_LIGHTS_OFF  = "C"
	STR_LIGHTS_ON   = "S"
	STR_CLEANER_OFF = "C"
	STR_CLEANER_ON  = "S"
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

	// fmt.Printf("%s\n", work_str)

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
	// filter:on aux2:on lights:off   | xxxTECD4S333333xxx
	// filter:on aux2:on lights:on    | xxxTESD4S333333xxx
	// filter:on aux2:off lights:on   | xxxTESD4C333333xxx
	// filter:off aux2:off lights:on  | xxxTDSD4C333333xxx
	// filter:off aux2:off lights:off | xxxTDCD4C333333xxx

	// TD..4.=filter-off, TE..4.=filter-on
	// T.C.4.=lights-off, T.S.4.=lights-on
	// T...4C=cleaner-off, T...4S=cleaner-on

	re = regexp.MustCompile("xxxT(.)(.)D4(.)......xxx")
	if len(re.FindStringSubmatch(work_str)) == 4 {

		switch re.FindStringSubmatch(work_str)[INDEX_FILTER] {

		case STR_FILTER_OFF:
			pool.FilterOn.Reading = 0
			pool.FilterOn.Last = time.Now()
		case STR_FILTER_ON:
			pool.FilterOn.Reading = 1
			pool.FilterOn.Last = time.Now()
		}

		switch re.FindStringSubmatch(work_str)[INDEX_LIGHTS] {

		case STR_LIGHTS_OFF:
			pool.LightOn.Reading = 0
			pool.LightOn.Last = time.Now()
		case STR_LIGHTS_ON:
			pool.LightOn.Reading = 1
			pool.LightOn.Last = time.Now()
		}

		switch re.FindStringSubmatch(work_str)[INDEX_CLEANER] {

		case STR_CLEANER_OFF:
			pool.CleanerOn.Reading = 0
			pool.CleanerOn.Last = time.Now()
		case STR_CLEANER_ON:
			pool.CleanerOn.Reading = 1
			pool.CleanerOn.Last = time.Now()
		}
	}
}
