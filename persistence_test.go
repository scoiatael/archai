package main_test

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/scoiatael/archai"
	. "github.com/scoiatael/archai/persistence"
	"github.com/scoiatael/archai/types"

	"github.com/gocql/gocql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testingKeyspace = "archai_test"

const findKeyspace = `select keyspace_name from system_schema.keyspaces where keyspace_name = ?`

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

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var _ = BeforeSuite(func() {
	var err error
	rand.Seed(time.Now().UnixNano())
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

	Describe("ReadEvents", func() {
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
	})
})
