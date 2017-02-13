package persistence

type Migration interface {
	Run(MigrationSession) error
}

type SimpleMigration struct {
	Query string
}

func (simpleMigration SimpleMigration) Run(session MigrationSession) error {
	err := session.Query(simpleMigration.Query)
	return err
}

var CreateEventsTable = SimpleMigration{Query: `
	CREATE TABLE IF NOT EXISTS events (
		id TIMEUUID,
		blob BLOB,
		stream VARCHAR,
		meta MAP<TEXT, TEXT>,
		PRIMARY KEY (stream, id)
	)
`}
