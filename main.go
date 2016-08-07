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
	"github.com/bccp-server/scheduler"
)

// Info from config file
type Config struct {
	Api_port       string
	Runner_service string
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
	mysql.Db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)
	go scheduler.DefaultScheduler.Start()
	go runners.WaitRunners(&scheduler.DefaultScheduler, config.Runner_service, config.Runner_token)
	api.SetupRestAPI(&wait, config.Api_port, config.Crt_file, config.Key_file)

	// Mysql tests
	mysql.Db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)

	wait.Wait()
}
