package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	if err := mainErr(); err != nil {
		log.Fatal(err)
	}
}
func mainErr() error {
	app := cli.NewApp()
	app.Name = "simulate parallel request to cassandra native transport"
	app.Commands = []*cli.Command{
		StressCommand(),
	}
	return app.Run(os.Args)
}

func StressCommand() *cli.Command {
	stressFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "seed",
			Aliases: []string{"s"},
			Value:   "localhost:9042",
			Usage:   "specify cassandra node",
		},
		&cli.IntFlag{
			Name:    "requests",
			Aliases: []string{"r"},
			Value:   1000,
			Usage:   "query count",
		},
		&cli.StringFlag{
			Name:     "mode",
			Aliases:  []string{"m"},
			Required: true,
			Usage:    "read/write",
		},
		&cli.IntFlag{
			Name:    "parallel",
			Aliases: []string{"p"},
			Value:   8,
			Usage:   "goroutine count",
		},
		&cli.IntFlag{
			Name:    "connection",
			Aliases: []string{"c"},
			Value:   8,
			Usage:   "connection per host",
		},
		&cli.DurationFlag{
			Name:    "timeout",
			Aliases: []string{"t"},
			Value:   5 * time.Second,
			Usage:   "connection timeout",
		},
		&cli.IntFlag{
			Name:  "cql",
			Value: 4,
			Usage: "cql version",
		},
		&cli.IntFlag{
			Name:  "replica-factor",
			Value: 1,
			Usage: "replica factor for keyspace",
		},
	}
	return &cli.Command{
		Name:   "stress",
		Usage:  "stress test for cassandra",
		Action: stress,
		Flags:  stressFlags,
	}
}
