package main_test

import (
	"fmt"

	. "github.com/scoiatael/archai"
	. "github.com/scoiatael/archai/persistence"
	"github.com/scoiatael/archai/types"

	"github.com/gocql/gocql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testingKeyspace = "archai_test"

const findKeyspace = `select keyspace_name from system.schema_keyspaces where keyspace_name = ?`

var (
	config    Config
	provider  CassandraProvider
	root_sess *gocql.Session
)

func testingKeyspaceExists() (bool, error) {
	var (
		err    error
		exists bool
	)

	iter := root_sess.Query(findKeyspace, testingKeyspace).Iter()
	exists = iter.Scan(nil)
	err = iter.Close()
	return exists, err
}

var dropKeyspace = fmt.Sprintf(`DROP KEYSPACE IF EXISTS %s`, testingKeyspace)

func dropTestingKeyspace() error {
	return root_sess.Query(dropKeyspace).Exec()
}

var dropMigrations = fmt.Sprintf(`DROP TABLE IF EXISTS %s.archai_migrations`, testingKeyspace)

func dropMigrationTable() error {
	return root_sess.Query(dropMigrations).Exec()
}

const findTable = `SELECT columnfamily_name from system.schema_columns where columnfamily_name = ? LIMIT 1 ALLOW FILTERING`

func migrationTableExists() (bool, error) {
	var (
		exists bool
		err    error
	)

	iter := root_sess.Query(findTable, "archai_migrations").Iter()
	exists = iter.Scan(nil)
	err = iter.Close()
	return exists, err
}

var _ = BeforeSuite(func() {
	var err error
	config = Config{Keyspace: testingKeyspace, Hosts: []string{"127.0.0.1"}}
	provider = CassandraProvider{Hosts: config.Hosts, Keyspace: config.Keyspace}
	cluster := provider.NewCluster()
	cluster.Consistency = gocql.All
	root_sess, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
	root_sess.Close()
})

var _ = Describe("Persistence", func() {
	Describe("MigrationSession", func() {
		BeforeEach(func() {
			err := dropTestingKeyspace()
			Expect(err).NotTo(HaveOccurred())
		})
		It("creates keyspace", func() {
			var (
				exists bool
				err    error
			)
			exists, err = testingKeyspaceExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())

			sess, err := provider.MigrationSession()
			Expect(err).NotTo(HaveOccurred())
			defer sess.Close()

			exists, err = testingKeyspaceExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

		})
	})
	Describe("ShouldRunMigration & DidRunMigration", func() {
		var (
			err  error
			sess MigrationSession
		)
		BeforeEach(func() {
			sess, err = provider.MigrationSession()
			Expect(err).NotTo(HaveOccurred())
			err = dropMigrationTable()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			sess.Close()
		})

		Context("when there's no migration table", func() {
			It("creates migration table", func() {
				var (
					exists bool
				)
				exists, err = migrationTableExists()
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())

				_, err := sess.ShouldRunMigration("foo")
				Expect(err).NotTo(HaveOccurred())

				exists, err = migrationTableExists()
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})

		Context("when migrations were not run", func() {
			It("returns true", func() {
				should, err := sess.ShouldRunMigration("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(should).To(BeTrue())
			})
		})

		Context("after migration was run", func() {
			It("returns false", func() {
				err := sess.DidRunMigration("foo")
				Expect(err).NotTo(HaveOccurred())

				should, err := sess.ShouldRunMigration("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(should).To(BeFalse())
			})
		})
	})

	Describe("ReadEvents", func() {
		var (
			sess Session
		)
		BeforeEach(func() {
			var (
				err error
			)
			err = dropTestingKeyspace()
			Expect(err).NotTo(HaveOccurred())
			s, err := provider.MigrationSession()
			Expect(err).NotTo(HaveOccurred())
			defer s.Close()
			err = CreateEventsTable.Run(s)
			Expect(err).NotTo(HaveOccurred())

			sess, err = provider.Session()
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			sess.Close()
		})
		Context("when there are no events", func() {
			It("returns empty array", func() {
				es, err := sess.ReadEvents("test-stream", "00000000-0000-1000-8080-808080808080", 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(es).To(BeEmpty())
			})
		})
		Context("after some events were added", func() {
			var (
				err    error
				events []types.Event
				cursor string
			)
			BeforeEach(func() {
				err = sess.WriteEvent("test-stream", []byte(`{ "a": 1 }`), make(map[string]string))
				Expect(err).NotTo(HaveOccurred())
				err = sess.WriteEvent("test-stream", []byte(`{ "a": 2 }`), make(map[string]string))
				Expect(err).NotTo(HaveOccurred())
				cursor = "00000000-0000-1000-8080-808080808080"
				events, err = sess.ReadEvents("test-stream", cursor, 10)
				Expect(err).NotTo(HaveOccurred())
			})
			JustBeforeEach(func() {
				events, err = sess.ReadEvents("test-stream", cursor, 10)
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns non-empty array", func() {
				Expect(events).NotTo(BeEmpty())
				Expect(events).To(HaveLen(2))
			})
			Context("when given cursor", func() {
				BeforeEach(func() {
					cursor = events[0].ID
				})
				It("returns events after cursor", func() {
					Expect(events).NotTo(BeEmpty())
					Expect(events).To(HaveLen(1))
				})
			})
		})
	})
})
