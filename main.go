package main

import (
	"os"

	"github.com/scoiatael/archai/actions"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "archai"
	app.Usage = "eventstore replacement"
	app.Version = Version
	app.Action = func(c *cli.Context) error {
		config := Config{Keyspace: "archai_test", Hosts: []string{"127.0.0.1"}}
		config.Actions = []actions.Action{
			//actions.WriteEventFromStream{Stream: "test-stream", Input: os.Stdin},
			actions.ReadEventsToStream{Stream: "test-stream", Output: os.Stdout},
			actions.HttpServer{Port: 8080},
		}
		return config.Run()
	}

	app.Run(os.Args)
}
