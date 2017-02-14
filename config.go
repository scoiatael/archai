package main

import (
	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/persistence"
)

type Config struct {
	keyspace string
	hosts    []string
	actions  []actions.Action
}

func (c Config) Migrations() map[string]persistence.Migration {
	m := make(map[string]persistence.Migration)
	m["create_events_table"] = persistence.CreateEventsTable
	return m
}

func (c Config) Persistence() persistence.Provider {
	hosts := make([]string, 1)
	hosts[0] = "127.0.0.1"
	provider := persistence.CassandraProvider{Hosts: hosts,
		Keyspace: c.keyspace}
	return &provider
}

func (c Config) Version() string {
	return Version
}
