package persistence

import (
	"fmt"

	"log"

	"github.com/gocql/gocql"
)

type MigrationSession interface {
	ShouldRunMigration(string) (bool, error)
	DidRunMigration(string) error
	Query(string) error
}

type CassandraMigrationSession struct {
	session *gocql.Session
}

func (sess *CassandraMigrationSession) Query(query string) error {
	return sess.session.Query(query).Exec()
}

const migrationTable = "archai_migrations"

var createMigrationTable = fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        name VARCHAR,
        PRIMARY KEY (name)
    )
`, migrationTable)

var findMigration = fmt.Sprintf(`SELECT name FROM %s WHERE name = ? LIMIT 1`, migrationTable)

var insertMigration = fmt.Sprintf(`INSERT INTO %s (name) VALUES (?)`, migrationTable)

func (sess *CassandraMigrationSession) ShouldRunMigration(name string) (bool, error) {
	if err := sess.session.Query(createMigrationTable).Exec(); err != nil {
		return false, err
	}
	log.Println("Looking for migration ", name)
	iter := sess.session.Query(findMigration, name).Iter()
	found := iter.Scan(nil)
	err := iter.Close()
	return !found, err
}

func (sess *CassandraMigrationSession) DidRunMigration(name string) error {
	return sess.session.Query(insertMigration, name).Exec()
}
