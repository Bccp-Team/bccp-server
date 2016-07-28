package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) AddRepoToNamespace(namespace string, repo string) (int, error) {
	req := "INSERT INTO namespace_repos VALUES(NULL,?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(namespace, repo)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return -1, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (db *Database) DeleteRepoFromNamespace(namespace string, repo string) error {
	req := "delete from namespace_repos where namespace=? and repo=?"
	del, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add runner: ", err.Error())
		return err
	}
	defer del.Close()

	_, err = del.Exec(namespace, repo)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return err
	}
	return nil
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
