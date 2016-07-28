package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/bccp-server/api"
	"github.com/bccp-server/mysql"
	"github.com/bccp-server/runners"
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
func ReadConfig(configfile string) Config {
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
	var configPath string

	flag.StringVar(&configPath, "config", "/etc/bccp/bccp.conf", "config path")
	flag.Parse()

	var config = ReadConfig(configPath)
	var wait sync.WaitGroup
	go runners.WaitRunners(config.Runner_port, config.Runner_token)
	api.SetupRestAPI(&wait, config.Port, config.Crt_file, config.Key_file)
	var db mysql.Database
	db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)
	for _, r := range db.ListRuns() {
		println("(", r.Id, " ", r.Status, " ", r.Runner_id, ")")
	}
	r := db.GetRun(1)
	println("(", r.Id, " ", r.Status, " ", r.Runner_id, ")")
	println(db.AddRunner())
	wait.Wait()
}
