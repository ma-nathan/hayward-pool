package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"strconv"
)

func handle_command_line_args() (heater_temp int) {

	heater_temp = POOL_TEMP_TARGET

	usage := `
Usage: pool [--heater_temp=<degrees_f>]

  Options:
    --heater_temp=<degrees_f>   Tell us the target pool temp setting.
    --version                   Show version.
    -h --help                   Show this screen.`

	arguments, err := docopt.Parse(usage, nil, true, "pool "+version, false)

	if err != nil {
		fmt.Printf("Error: \"%v\"\n", err)
		os.Exit(1)
	}

	if arguments["--heater_temp"] != nil && arguments["--heater_temp"] != "" {
		heater_temp, _ = strconv.Atoi(arguments["--heater_temp"].(string))
	}

	return
}
