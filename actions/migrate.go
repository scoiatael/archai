package actions

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// Migrate implements Action interface
type Migrate struct{}

// Run all migrations
func (a Migrate) Run(c Context) error {
	migrationSession, err := c.Persistence().MigrationSession()
	if err != nil {
		return errors.Wrap(err, "Obtaining session failed")
	}
	defer migrationSession.Close()
	for name, m := range c.Migrations() {
		shouldRun, err := migrationSession.ShouldRunMigration(name)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Checking if should run migration %s failed", name))
		}
		if shouldRun {
			log.Println("Executing ", name)
			err = m.Run(migrationSession)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Running migration %s failed", name))
			}
			err = migrationSession.DidRunMigration(name)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Running after-migration hook for %s failed", name))
			}
		}
	}
	return nil
}

func (a Migrate) MarshalJSON() ([]byte, error) {
	return []byte(`"Migrate Cassandra keyspace"`), nil
}
