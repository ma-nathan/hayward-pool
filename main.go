package main

import (
	"fmt"
)

//const (
//   POOL_TEMP_TARGET_INVALID = -1
//)

func main() {

	fmt.Println("pool-data-collector polls a Hayward Aqua Connect Local network device.")

	//    target_temp := handle_command_line_args()
	var config = ReadConfig()

	//	if target_temp == POOL_TEMP_TARGET_INVALID {
	//		target_temp=config.PoolTempTarget
	//	}

	c := influxDBClient(config)

	go update_datastore(c, config)
	watch_http_endpoint(config)
}
