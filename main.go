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
		action := actions.ReadEventsToStream{Stream: "testing-stream", Output: os.Stdout}
		err := action.Run(config)
		return err
	}

	app.Run(os.Args)
}
