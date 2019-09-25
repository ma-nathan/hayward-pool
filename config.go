package main

import (
	"os"
	"log"
	"github.com/BurntSushi/toml"
)


// $ cat settings.ini
//
// PoolHost="pool"
// DatabaseURL="http://metrics:8086"
// DatabaseUser="admin"
// DatabasePassword="blahblah"
// DatabaseDatabase="scrape"
//

type Config struct {
	PoolHost			string
	PoolTempTarget		int
	DatabaseURL			string
	DatabaseUser		string
	DatabasePassword	string
	DatabaseDatabase	string
}

var configfile = "settings.ini"


func ReadConfig() Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}


