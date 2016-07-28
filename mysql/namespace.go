package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// List Runs
func (db *Database) ListNamespaces() ([]string, error) {

	var name string
	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM run")
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	var namespaces []string

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&name)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		namespaces = append(namespaces, name)
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return namespaces, err
	}

	return namespaces, nil
}

func (db *Database) AddRepoToNamepace(namespace string, depo string) (int, error) {
	req := "INSERT INTO namespace_repos VALUES(NULL,?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(namespace, depo)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return -1, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// Get namespace's repos
func (db *Database) GetNamespaceRepos(name string) ([]string, error) {

	var repo string
	// Execute the query
	req := "SELECT repo FROM namespace_repos where namespace='" + name + "'"
	rows, err := db.conn.Query(req)
	if err != nil {
		log.Print("ERROR: Unable to select namespace repos: ", err.Error())
		return nil, err
	}

	var namespace []string

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&repo)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		namespace = append(namespace, repo)
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return namespace, err
	}

	return namespace, nil
}
