package app

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// ExecApp - Configures and executes the Dohyo CLI application
func ExecApp() {
	app := &cli.App{
		Name:  "Dohyo",
		Usage: "Backend for Basho, the Fantasy Sumo App",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config",
				Usage: "name of the config file to use",
			},
			&cli.StringFlag{
				Name:  "config-dir",
				Value: ".",
				Usage: "directory where the config file resides",
			},
			&cli.StringFlag{
				Name:  "config-ext",
				Value: "yaml",
				Usage: "file type extension of the config file",
			},
			&cli.StringFlag{
				Name:  "bs-log",
				Value: "./logs/bootstrap.log",
				Usage: "file for writing the bootstrap logs",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "error",
				Usage: "logging level (Error, Warn, Info, Debug)",
			},
			&cli.BoolFlag{
				Name:  "log-clear",
				Usage: "clears the log files before execution",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run the server",
				Action: run,
			},
			{
				Name:   "activate",
				Usage:  "sets the active tournament",
				Action: activate,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "tournament name to set as the active tournament",
						Required: true,
					},
				},
			},
			{
				Name:   "load",
				Usage:  "loads data into the database from data files",
				Action: load,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "data-file",
						Usage:    "data file to source the data from",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "clear",
						Usage: "clear the tables before adding data",
					},
				},
			},
			{
				Name:   "migrate",
				Usage:  "migrates the provided tables into the database",
				Action: migrate,
			},
		},
	}

	fmt.Println("running Dohyo app")
	if err := app.Run(os.Args); err != nil {
		panic(fmt.Errorf("fatal error running app: %s", err))
	}
}
