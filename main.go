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
		config := Config{keyspace: "archai_test"}
		action := actions.ReadEvent{Stream: "testing-stream", Amount: 10}
		err := action.Run(config)
		return err
	}

	app.Run(os.Args)
}
