package main

import (
	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/persistence"
)

// Config is a context for all application actions.
type Config struct {
	Keyspace string
	Hosts    []string
	Actions  []actions.Action
}

func (c Config) Migrations() map[string]persistence.Migration {
	m := make(map[string]persistence.Migration)
	m["create_events_table"] = persistence.CreateEventsTable
	return m
}

func (c Config) Persistence() persistence.Provider {
	provider := persistence.CassandraProvider{Hosts: c.Hosts, Keyspace: c.Keyspace}
	return &provider
}

// Version returns current version
func (c Config) Version() string {
	return Version
}
