package main

import (
	"fmt"
	"log"

	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/persistence"
)

// Config is a context for all application actions.
type Config struct {
	Keyspace    string
	Hosts       []string
	Actions     []actions.Action
	provider    persistence.Provider
	initialized bool
}

func (c Config) HandleErr(err error) {
	log.Print(err)
}

func (c Config) Migrations() map[string]persistence.Migration {
	m := make(map[string]persistence.Migration)
	m["001_create_events_table"] = persistence.CreateEventsTable
	return m
}

func (c Config) Persistence() persistence.Provider {
	if !c.initialized {
		panic(fmt.Errorf("Persistence not initialized!"))
	}
	return c.provider
}

// Version returns current version
func (c Config) Version() string {
	return Version
}

func (c Config) HttpHandler() actions.HttpHandler {
	return http.NewIris(c)
}

func (c *Config) Init() error {
	new_provider := persistence.CassandraProvider{Hosts: c.Hosts, Keyspace: c.Keyspace}
	err := new_provider.Init()
	if err != nil {
		return err
	}
	c.provider = &new_provider
	c.initialized = true
	return nil
}

func (c Config) Run() error {
	if err := c.Init(); err != nil {
		return err
	}
	for _, a := range c.Actions {
		err := a.Run(c)
		if err != nil {
			return err
		}
	}
	return nil
}
