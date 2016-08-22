package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Bccp-Team/bccp-server/api"
	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"
	"github.com/Bccp-Team/bccp-server/scheduler"
	"github.com/BurntSushi/toml"
)

type Config struct {
	APIPort       string
	RunnerService string
	RunnerToken   string
	KeyFile       string
	CrtFile       string
	MysqlDatabase string
	MysqlUser     string
	MysqlPassword string
}

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
	mysql.Db.Connect(config.MysqlDatabase, config.MysqlUser, config.MysqlPassword)

	go scheduler.DefaultScheduler.Start()
	go runners.WaitRunners(&scheduler.DefaultScheduler, config.RunnerService, config.RunnerToken)

	api.SetupRestAPI(&wait, config.APIPort, config.CrtFile, config.KeyFile)

	// Mysql tests
	mysql.Db.Connect(config.MysqlDatabase, config.MysqlUser, config.MysqlPassword)

	wait.Wait()
}
