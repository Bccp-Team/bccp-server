package mysql

import (
	"database/sql"
	"log"

	. "github.com/Bccp-Team/bccp-server/proto/api"
)

func (db *Database) AddRepoToNamespace(namespace string, repo string, ssh string) (int, error) {
	req := "INSERT INTO namespace_repos VALUES(NULL,?,?,?)"
	insert, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add namespace_repo: ", err.Error())
		return -1, err
	}
	defer insert.Close()

	res, err := insert.Exec(namespace, repo, ssh)
	if err != nil {
		log.Print("ERROR: Unable to insert namespace_repo: ", err.Error())
		return -1, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (db *Database) DeleteRepoFromNamespace(namespace string, repo int64) error {
	req := "delete from namespace_repos where namespace=? and repo=?"
	del, err := db.conn.Prepare(req)
	if err != nil {
		log.Print("ERROR: Unable to prepare add namespace_repo: ", err.Error())
		return err
	}
	defer del.Close()

	_, err = del.Exec(namespace, repo)
	if err != nil {
		log.Print("ERROR: Unable to delete namespace_repo: ", err.Error())
		return err
	}
	return nil
}

// Get namespace's repos
func (db *Database) GetNamespaceRepos(name *string) ([]*Repo, error) {
	var repo string
	var ssh string
	var id int64
	var namespace string

	var rows *sql.Rows
	var err error

	// Execute the query
	if name == nil {
		req := "SELECT id, repo, ssh, namespace FROM namespace_repos"
		rows, err = db.conn.Query(req)
	} else {
		req := "SELECT id, repo, ssh, namespace FROM namespace_repos where namespace=?"
		rows, err = db.conn.Query(req, *name)
	}
	if err != nil {
		log.Print("ERROR: Unable to select namespace repos: ", err.Error())
		return nil, err
	}

	var repos []*Repo

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&id, &repo, &ssh, &namespace)
		if err != nil {
			log.Print("ERROR: Unable to get next row: ", err.Error())
			return nil, err
		}

		repos = append(repos, &Repo{Id: id, Repo: repo, Ssh: ssh, Namespace: namespace})
	}
	if err = rows.Err(); err != nil {
		log.Print("ERROR: Undefined row err: ", err.Error())
		return nil, err
	}

	return repos, nil
}

func (db *Database) GetRepo(id int64) (*Repo, error) {
	var repo string
	var ssh string
	var namespace string

	// Execute the query
	req := "SELECT repo, ssh, namespace FROM namespace_repos where id=?"
	err := db.conn.QueryRow(req, id).Scan(&repo, &ssh, &namespace)
	if err != nil {
		log.Print("ERROR: Unable to select namespace repos: ", err.Error())
		return nil, err
	}

	return &Repo{Id: id, Repo: repo, Ssh: ssh, Namespace: namespace}, nil
}

func (db *Database) GetRepoFromName(name string, namespace string) (*Repo, error) {
	var id int64
	var repo string
	var ssh string

	// Execute the query
	req := "SELECT id, repo, ssh FROM namespace_repos where repo=? AND namespace=?"
	err := db.conn.QueryRow(req, name, namespace).Scan(&id, &repo, &ssh)
	if err != nil {
		log.Print("ERROR: Unable to select namespace repos: ", err.Error())
		return nil, err
	}

	return &Repo{Id: id, Repo: repo, Ssh: ssh, Namespace: namespace}, nil
}
