package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Bccp-Team/bccp-server/api/rest"
	"github.com/Bccp-Team/bccp-server/api/rpc"
	"github.com/Bccp-Team/bccp-server/mysql"
	"github.com/Bccp-Team/bccp-server/runners"
	"github.com/Bccp-Team/bccp-server/scheduler"
	"github.com/BurntSushi/toml"
)

type Config struct {
	APIPort       string
	RPCPort       string
	RunnerService string
	RunnerToken   string
	KeyFile       string
	CrtFile       string
	MysqlService  string
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
	var wait sync.WaitGroup

	flag.StringVar(&configPath, "config", "/etc/bccp/bccp.conf", "config path")
	flag.Parse()

	var config = ReadConfig(configPath)
	mysql.Db.Connect(config.MysqlService, config.MysqlDatabase, config.MysqlUser, config.MysqlPassword)

	go scheduler.DefaultScheduler.Start()
	go runners.WaitRunners(&scheduler.DefaultScheduler, config.RunnerService, config.RunnerToken)

	rest.SetupRestAPI(&wait, config.APIPort, config.CrtFile, config.KeyFile)

	//wait.Wait()
	err := rpc.SetupRpc(config.RPCPort)

	log.Panic(err)
}
