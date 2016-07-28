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

	// Mysql tests
	mysql.Db.Connect(config.Mysql_database, config.Mysql_user, config.Mysql_password)
	id, _ := mysql.Db.AddRun("test")
	mysql.Db.LaunchRun(id, 1)
	mysql.Db.UpdateRunStatus(id, "timeout")
	mysql.Db.UpdateRunLogs(id, "lol")
	mysql.Db.UpdateRunLogs(id, "lel")
	mysql.Db.UpdateRunLogs(id, "lol")
	id, _ = mysql.Db.AddRun("bite")
	id, _ = mysql.Db.AddRun("poil")
	runs, _ := mysql.Db.ListRuns()
	for _, run := range runs {
		println("(", run.Id, run.Status, run.Runner_id, run.Repo, run.Logs, ")")
	}

	wait.Wait()
}
