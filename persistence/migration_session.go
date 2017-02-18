package persistence

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

type MigrationSession interface {
	ShouldRunMigration(string) (bool, error)
	DidRunMigration(string) error
	Exec(string) error
	Close()
}

type CassandraMigrationSession struct {
	*gocql.Session
}

func (sess *CassandraMigrationSession) Exec(query string) error {
	return sess.Query(query).Exec()
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

func (sess *CassandraMigrationSession) createMigrationTableIfNeeded() error {
	if err := sess.Query(createMigrationTable).Exec(); err != nil {
		return errors.Wrap(err, "Query to createMigrationTable failed")
	}
	return nil
}

func (sess *CassandraMigrationSession) ShouldRunMigration(name string) (bool, error) {
	if err := sess.createMigrationTableIfNeeded(); err != nil {
		return false, errors.Wrap(err, "Failed to create migrations table if needed")
	}
	iter := sess.Query(findMigration, name).Iter()
	found := iter.Scan(nil)
	err := iter.Close()
	return !found, errors.Wrap(err, "Closing iterator for findMigration failed")
}

func (sess *CassandraMigrationSession) DidRunMigration(name string) error {
	if err := sess.createMigrationTableIfNeeded(); err != nil {
		return err
	}
	return sess.Query(insertMigration, name).Exec()
}
