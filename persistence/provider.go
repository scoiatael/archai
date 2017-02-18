package persistence

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

type Provider interface {
	Session() (Session, error)
	MigrationSession() (MigrationSession, error)
}

type CassandraProvider struct {
	Hosts    []string
	Keyspace string
}

func (cp *CassandraProvider) newCluster() *gocql.ClusterConfig {
	return gocql.NewCluster(cp.Hosts...)
}

func (cp *CassandraProvider) Session() (Session, error) {
	cluster := cp.newCluster()
	cluster.Keyspace = cp.Keyspace
	cluster.Consistency = gocql.Quorum
	sess, err := cluster.CreateSession()
	return &CassandraSession{sess}, errors.Wrap(err, "CreateSession failed")
}

const createKeySpace = `CREATE KEYSPACE IF NOT EXISTS %s WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`

func (cp *CassandraProvider) MigrationSession() (MigrationSession, error) {
	cluster := cp.newCluster()
	cluster.Consistency = gocql.All
	sess, err := cluster.CreateSession()
	if err != nil {
		return &CassandraMigrationSession{}, errors.Wrap(err, "CreateSession failed")
	}
	err = sess.Query(fmt.Sprintf(createKeySpace, cp.Keyspace)).Exec()
	if err != nil {
		return &CassandraMigrationSession{}, errors.Wrap(err, "Query to CreateKeyspace failed")
	}
	cluster.Keyspace = cp.Keyspace
	sess, err = cluster.CreateSession()

	return &CassandraMigrationSession{sess}, errors.Wrap(err, "CreateSession failed")
}
