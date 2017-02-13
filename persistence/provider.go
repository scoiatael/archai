package persistence

import (
	"fmt"

	"github.com/gocql/gocql"
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
	return &CassandraSession{session: sess}, err
}

func (cp *CassandraProvider) MigrationSession() (MigrationSession, error) {
	cluster := cp.newCluster()
	cluster.Consistency = gocql.All
	sess, err := cluster.CreateSession()
	if err != nil {
		return &CassandraMigrationSession{}, err
	}
	err = sess.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`, cp.Keyspace)).Exec()
	if err != nil {
		return &CassandraMigrationSession{}, err
	}
	cluster.Keyspace = cp.Keyspace
	sess, err = cluster.CreateSession()

	return &CassandraMigrationSession{session: sess}, err
}
