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
		return actions.Migrate{}.Run(config)
	}

	app.Run(os.Args)
}
