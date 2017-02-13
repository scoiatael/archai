package actions

import (
	"log"
)

// Migrate implements Action interface
type Migrate struct{}

// Run all migrations
func (a Migrate) Run(c Context) error {
	persistenceProvider := c.Persistence()
	migrationSession, err := persistenceProvider.MigrationSession()
	if err != nil {
		return err
	}
	for name, m := range c.Migrations() {
		shouldRun, err := migrationSession.ShouldRunMigration(name)
		if err != nil {
			return err
		}
		if shouldRun {
			log.Println("Executing ", name)
			err = m.Run(migrationSession)
			if err != nil {
				return err
			}
			err = migrationSession.DidRunMigration(name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
