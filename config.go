package main

import (
	"fmt"
	"log"

	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/http"
	"github.com/scoiatael/archai/persistence"
	"github.com/scoiatael/archai/telemetry"
)

// Config is a context for all application actions.
type Config struct {
	Keyspace   string
	Hosts      []string
	Actions    []actions.Action
	StatsdAddr string
	Features   map[string]bool

	provider    persistence.Provider
	telemetry   telemetry.Datadog
	initialized bool
}

func (c Config) HandleErr(err error) {
	log.Print(err)
	c.Telemetry().Failure("server_error", err.Error())
}

func (c Config) Migrations() map[string]persistence.Migration {
	m := make(map[string]persistence.Migration)
	m["001_create_events_table"] = persistence.CreateEventsTable
	return m
}

func (c Config) Persistence() persistence.Provider {
	if !c.initialized {
		panic(fmt.Errorf("Config not initialized!"))
	}
	return c.provider
}

func (c Config) Telemetry() telemetry.Datadog {
	if !c.initialized {
		panic(fmt.Errorf("Config not initialized!"))
	}
	return c.telemetry
}

// Version returns current version
func (c Config) Version() string {
	return Version
}

func (c Config) HttpHandler() actions.HttpHandler {
	return http.NewIris(c, c.Features["dev_logger"])
}

func (c *Config) Init() error {
	new_provider := persistence.CassandraProvider{Hosts: c.Hosts, Keyspace: c.Keyspace}
	err := new_provider.Init()
	if err != nil {
		return err
	}
	c.provider = &new_provider

	dd := telemetry.NewDatadog(c.StatsdAddr, c.Keyspace)
	c.telemetry = dd

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
