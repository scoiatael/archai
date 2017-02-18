package main

import (
	"log"

	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/persistence"
)

// Config is a context for all application actions.
type Config struct {
	Keyspace string
	Hosts    []string
	Actions  []actions.Action
}

func (c Config) HandleErr(err error) {
	log.Print(err)
	panic(err)
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

func (c Config) HttpHandler() actions.HttpHandler {
	return &http.FastHttpHandler{Context: c}
}

func (c Config) Run() error {
	for _, a := range c.Actions {
		err := a.Run(c)
		if err != nil {
			return err
		}
	}
	return nil
}
