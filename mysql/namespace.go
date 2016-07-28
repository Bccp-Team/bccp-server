package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// List Runs
func (db *Database) ListNamespaces() ([]string, error) {

	var name string
	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM namespace")
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

func (db *Database) AddNamespace(namespace string) error {
	req := "INSERT INTO namespace VALUES(?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add namespace: ", err.Error())
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(namespace)
	if err != nil {
		log.Print("ERROR: Unable to insert runner: ", err.Error())
		return err
	}
	return nil
}

func (db *Database) DeleteNamespace(namespace string) error {
	repos, err := db.GetNamespaceRepos(namespace)
	if err != nil {
		log.Print("ERROR: Unable to prepare delete namespace: ", err.Error())
		return err
	}

	for _, repo := range repos {
		db.DeleteRepoFromNamespace(namespace, repo)
	}

	req := "delete from namespace where name=?"
	del, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare delete namespace: ", err.Error())
		return err
	}
	defer del.Close()

	_, err = del.Exec(namespace)
	if err != nil {
		log.Print("ERROR: Unable to delete namespace: ", err.Error())
		return err
	}

	return nil
}
