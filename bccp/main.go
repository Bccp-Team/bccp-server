package main

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/bccp/api"
)

// Info from config file
type Config struct {
	server_key_file string
	server_crt_file string
	bccp_database   string
	bccp_user       string
	bccp_password   string
}

// Reads info from config file
func ReadConfig() Config {
	var configfile = "/etc/bccp/bccp.conf"
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

func main() {
	//var config = ReadConfig()
	var wait sync.WaitGroup
	api.SetupRestAPI(&wait)
	wait.Wait()
}
