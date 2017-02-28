package persistence_test

import (
	"fmt"

	. "github.com/scoiatael/archai/persistence"
	"github.com/scoiatael/archai/types"
	"github.com/scoiatael/archai/util"

	"github.com/gocql/gocql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testingKeyspace = "archai_test"

const findKeyspace = `select keyspace_name from system_schema.keyspaces where keyspace_name = ?`

var (
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

const findTable = `SELECT table_name from system_schema.columns where keyspace_name = ? AND table_name = ? LIMIT 1 ALLOW FILTERING`

func migrationTableExists() (bool, error) {
	var (
		exists bool
		err    error
	)

	iter := root_sess.Query(findTable, testingKeyspace, "archai_migrations").Iter()
	exists = iter.Scan(nil)
	err = iter.Close()
	return exists, err
}

func randomString() string {
	return util.RandomString(10)
}

var _ = BeforeSuite(func() {
	var err error
	provider = CassandraProvider{Hosts: []string{"127.0.0.1"}, Keyspace: testingKeyspace}
	cluster := provider.NewCluster()
	cluster.Consistency = gocql.All
	root_sess, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {
	root_sess.Query(fmt.Sprintf("truncate table %s.events", testingKeyspace)).Exec()
	root_sess.Close()
})

var _ = Describe("Persistence", func() {
	Describe("MigrationSession", func() {
		It("creates keyspace", func() {
			var (
				exists bool
				err    error
			)
			exists, err = testingKeyspaceExists()
			Expect(err).NotTo(HaveOccurred())
			if exists {
				Skip("Testing keyspace already exists; drop it to run this test")
			}

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
		})

		// AfterEach(func() {
		// 	sess.Close()
		// })

		Context("when there's no migration table", func() {
			It("creates migration table", func() {
				var (
					exists bool
				)
				exists, err = migrationTableExists()
				Expect(err).NotTo(HaveOccurred())
				if exists {
					Skip("Migration table exists; drop it to run this test")
				}

				_, err := sess.ShouldRunMigration(randomString())
				Expect(err).NotTo(HaveOccurred())

				exists, err = migrationTableExists()
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})

		Context("when migrations were not run", func() {
			It("returns true", func() {
				should, err := sess.ShouldRunMigration(randomString())
				Expect(err).NotTo(HaveOccurred())
				Expect(should).To(BeTrue())
			})
		})

		Context("after migration was run", func() {
			It("returns false", func() {
				name := randomString()
				err := sess.DidRunMigration(name)
				Expect(err).NotTo(HaveOccurred())

				should, err := sess.ShouldRunMigration(name)
				Expect(err).NotTo(HaveOccurred())
				Expect(should).To(BeFalse())
			})
		})
	})

	Describe("ReadEvents & WriteEvent", func() {
		var (
			sess Session
		)
		BeforeEach(func() {
			var (
				err error
			)
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
				es, err := sess.ReadEvents(randomString(), "00000000-0000-1000-8080-808080808080", 10)
				Expect(err).NotTo(HaveOccurred())
				Expect(es).To(BeEmpty())
			})
		})
		Context("after some events were added", func() {
			var (
				err    error
				events []types.Event
				cursor string
				stream string
			)
			BeforeEach(func() {
				stream = randomString()
				err = sess.WriteEvent(stream, []byte(`{ "a": 1 }`), make(map[string]string))
				Expect(err).NotTo(HaveOccurred())
				err = sess.WriteEvent(stream, []byte(`{ "a": 2 }`), make(map[string]string))
				Expect(err).NotTo(HaveOccurred())
				cursor = "00000000-0000-1000-8080-808080808080"
				events, err = sess.ReadEvents(stream, cursor, 10)
				Expect(err).NotTo(HaveOccurred())
			})
			JustBeforeEach(func() {
				events, err = sess.ReadEvents(stream, cursor, 10)
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
		Describe("WriteEvent", func() {
			Measure("small events", func(b Benchmarker) {
				var (
					err    error
					stream string
				)
				b.Time("", func() {
					stream = randomString()
					err = sess.WriteEvent(stream, []byte(`{ "a": 1 }`), make(map[string]string))
					Expect(err).NotTo(HaveOccurred())
				})
			}, 20)

			Measure("big events", func(b Benchmarker) {
				var (
					err    error
					stream string
					blob   []byte
				)
				blob = make([]byte, 1024*1024)
				b.Time("", func() {
					stream = randomString()
					err = sess.WriteEvent(stream, blob, make(map[string]string))
					Expect(err).NotTo(HaveOccurred())
				})
			}, 1)
		})
	})
})
