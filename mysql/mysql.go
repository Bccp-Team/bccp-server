package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	conn *sql.DB
}

func (db *Database) Connect(database string, user string, password string) {

	dsn := user + ":" + password + "@/" + database + "?parseTime=true"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("ERROR: Unable to open mysql connection: ", err.Error())
	}
	db.conn = conn

	err = db.conn.Ping()
	if err != nil {
		log.Fatal("ERROR: Unable to ping database: ", err.Error())
	}
}
