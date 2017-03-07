package main

import (
	"os"
	"strings"

	"github.com/scoiatael/archai/actions"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "archai"
	app.Usage = "eventstore replacement"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "migrate",
			Usage: "Migrate Cassandra on startup?",
		},
		cli.BoolFlag{
			Name:  "dev-logger",
			Usage: "Enable dev logger",
		},
		cli.StringFlag{
			Name:  "keyspace",
			Value: "archai",
			Usage: "Cassandra keyspace to operate in",
		},
		cli.StringFlag{
			Name:  "hosts",
			Value: "127.0.0.1",
			Usage: "Comma-separated list of Cassandra hosts",
		},
		cli.StringFlag{
			Name:  "listen",
			Value: "127.0.0.1",
			Usage: "Address to listen on",
		},
		cli.StringFlag{
			Name:  "telemetry",
			Value: "127.0.0.1",
			Usage: "Address to send metrics to",
		},
		cli.Int64Flag{
			Name:  "port",
			Value: 8080,
			Usage: "Port to listen on",
		},
	}
	app.Action = func(c *cli.Context) error {
		config := Config{}
		config.Features = make(map[string]bool)

		config.Keyspace = c.String("keyspace")
		config.Hosts = strings.Split(c.String("hosts"), ",")
		config.StatsdAddr = c.String("telemetry")
		if c.Bool("migrate") {
			config.Append(actions.Migrate{})
		}
		if c.Bool("dev-logger") {
			config.Features["dev_logger"] = true
		}
		config.Append(actions.HttpServer{
			Port: c.Int("port"),
			Addr: c.String("listen")})
		config.PrettyPrint()
		return config.Run()
	}

	app.Run(os.Args)
}
