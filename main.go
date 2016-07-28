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
	mysql.Db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)
	go scheduler.DefaultScheduler.Start()
	go runners.WaitRunners(&scheduler.DefaultScheduler, config.Runner_port, config.Runner_token)
	api.SetupRestAPI(&wait, config.Port, config.Crt_file, config.Key_file)

	id, _ := mysql.Db.AddRun("foobarbaz")

	scheduler.DefaultScheduler.AddRun(id)

	// Mysql tests
	mysql.Db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)

	mysql.Db.AddNamespace("tc1")
	mysql.Db.AddNamespace("tc2")
	mysql.Db.AddBatch("tc1", "#!", 10, 120)
	mysql.Db.AddBatch("tc2", "#!", 10, 120)

	println("-----------")
	ns, _ := mysql.Db.ListNamespaces()
	for _, n := range ns {
		println("(", n, ")")
	}
	println("-----------")
	batchs := mysql.Db.ListBatchs()
	for _, n := range batchs {
		println("(", n.Id, n.Namespace, n.Init_script, n.Update_time, n.Timeout, ")")
	}
	println("-----------")

	wait.Wait()
}
