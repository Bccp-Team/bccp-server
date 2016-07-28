package main

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/bccp/api"
	"github.com/bccp/mysql"
)

// Info from config file
type Config struct {
	Port           string
	Runner_port    string
	Runner_token   string
	Key_file       string
	Crt_file       string
	Mysql_database string
	Mysql_user     string
	Mysql_password string
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
	return config
}

func main() {
	var config = ReadConfig()
	var wait sync.WaitGroup
	go WaitRunners(config.Runner_port, config.Runner_token)
	api.SetupRestAPI(&wait, config.Port, config.Crt_file, config.Key_file)
	var db mysql.Database
	db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)
	for _, r := range db.ListRunners() {
		println("(", r.Id, " ", r.Status, ")")
	}
	r := db.GetRunner(1)
	println("(", r.Id, " ", r.Status, ")")
	wait.Wait()
}
