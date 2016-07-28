package main

import (
	"flag"
	"fmt"
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
	for _, r := range db.ListRunners() {
		print("(", r.Id, ", ", r.Status, ", ")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d )\n",
			r.Last_connection.Year(), r.Last_connection.Month(), r.Last_connection.Day(),
			r.Last_connection.Hour(), r.Last_connection.Minute(), r.Last_connection.Second())
	}
	r := db.GetRunner(1)
	print("(", r.Id, ", ", r.Status, ", ")
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d )\n",
		r.Last_connection.Year(), r.Last_connection.Month(), r.Last_connection.Day(),
		r.Last_connection.Hour(), r.Last_connection.Minute(), r.Last_connection.Second())

	id := db.AddRunner("91.121.83.195")
	println(id)
	println(db.UpdateRunner(id, "running"))
	for _, r := range db.ListRunners() {
		print("(", r.Id, ", ", r.Status, ", ")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d )\n",
			r.Last_connection.Year(), r.Last_connection.Month(), r.Last_connection.Day(),
			r.Last_connection.Hour(), r.Last_connection.Minute(), r.Last_connection.Second())
	}
	wait.Wait()
}
