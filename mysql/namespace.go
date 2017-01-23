package mysql

import (
	"log"

	. "github.com/Bccp-Team/bccp-server/proto/api"
)

func (db *Database) ListNamespaces() ([]*Namespace, error) {

	var name string
	var is_ci bool
	// Execute the query
	rows, err := db.conn.Query("SELECT * FROM namespace")
	if err != nil {
		log.Print("ERROR: Unable to select run: ", err.Error())
		return nil, err
	}

	var namespaces []*Namespace

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&name, &is_ci)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		namespaces = append(namespaces, &Namespace{Name: name, IsCi: is_ci})
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return namespaces, err
	}

	return namespaces, nil
}

func (db *Database) AddNamespace(namespace string, is_ci bool) error {
	req := "INSERT INTO namespace VALUES(?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add namespace: ", err.Error())
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(namespace, is_ci)
	if err != nil {
		log.Print("ERROR: Unable to insert namespace: ", err.Error())
		return err
	}
	return nil
}

func (db *Database) GetNamespace(namespace string) (*Namespace, error) {
	var is_ci bool
	req := "SELECT is_ci FROM namespace where name=?"
	err := db.conn.QueryRow(req, namespace).Scan(&is_ci)
	if err != nil {
		log.Print("ERROR: Unable to get namespace: ", err.Error())
		return nil, err
	}
	return &Namespace{Name: namespace, IsCi: is_ci}, nil
}

func (db *Database) DeleteNamespace(namespace string) error {
	repos, err := db.GetNamespaceRepos(&namespace)
	if err != nil {
		log.Print("ERROR: Unable to prepare delete namespace: ", err.Error())
		return err
	}

	for _, repo := range repos {
		db.DeleteRepoFromNamespace(namespace, repo.Id)
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
