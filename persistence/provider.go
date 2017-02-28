package persistence

import (
	"fmt"
	"time"

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
	session  *Session
}

func (cp *CassandraProvider) NewCluster() *gocql.ClusterConfig {
	return gocql.NewCluster(cp.Hosts...)
}

func (cp *CassandraProvider) newSession() (Session, error) {
	cluster := cp.NewCluster()
	cluster.Keyspace = cp.Keyspace
	cluster.Consistency = gocql.Quorum
	sess, err := cluster.CreateSession()
	return &CassandraSession{sess}, errors.Wrap(err, "CreateSession failed")
}

func (cp *CassandraProvider) Session() (Session, error) {
	if cp.session != nil {
		return *cp.session, nil
	}
	return nil, fmt.Errorf("Initialize CassandraProvider with NewProvider first")
}

const createKeySpace = `CREATE KEYSPACE IF NOT EXISTS %s WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`

func (cp *CassandraProvider) MigrationSession() (MigrationSession, error) {
	cluster := cp.NewCluster()
	cluster.Timeout = 2000 * time.Millisecond
	cluster.Consistency = gocql.All
	sess, err := cluster.CreateSession()
	if err != nil {
		return &CassandraMigrationSession{}, errors.Wrap(err, "CreateSession failed")
	}
	defer sess.Close()
	err = sess.Query(fmt.Sprintf(createKeySpace, cp.Keyspace)).Exec()
	if err != nil {
		return &CassandraMigrationSession{}, errors.Wrap(err, "Query to CreateKeyspace failed")
	}
	cluster.Keyspace = cp.Keyspace
	sess, err = cluster.CreateSession()

	return &CassandraMigrationSession{sess}, errors.Wrap(err, "CreateSession failed")
}

func (c *CassandraProvider) Init() error {
	new_sess, err := c.newSession()
	if err == nil {
		c.session = &new_sess
	}
	return err
}
