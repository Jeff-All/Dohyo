// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Jeff-All/Dohyo/handlers"
	"github.com/Jeff-All/Dohyo/helpers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var bslog = logrus.New()
var log = logrus.New()
var routeHandlers = make(map[string]handlers.HandlerInterface)

func main() {
	fmt.Println("starting Dohyo")

	app := &cli.App{
		Name:   "Dohyo",
		Usage:  "Backend for Basho, the Fantasy Sumo App",
		Action: run,
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
		},
	}

	fmt.Println("running Dohyo app")
	if err := app.Run(os.Args); err != nil {
		panic(fmt.Errorf("fatal error running app: %s", err))
	}
}

func run(c *cli.Context) error {
	if err := loadBootstrapLog(c); err != nil {
		panic(fmt.Errorf("error loading bootstrap log: %s", err))
	}
	if err := loadConfig(c); err != nil {
		bslog.Errorf("error loading config: %s", err)
	}
	if err := loadLog(c); err != nil {
		bslog.Errorf("error loading log: %s", err)
	}
	buildHandlers()
	router := buildRouter()
	server := buildServer(router)
	bslog.Info("launching server")
	log.Fatal(server.ListenAndServe())
	return nil
}

func loadBootstrapLog(c *cli.Context) error {
	var err error
	fmt.Printf("loading bootstrap log: '%s'\n", c.String("bs-log"))
	if bslog.Out, err = helpers.AppendOrCreateFile(c.String("bs-log")); err != nil {
		return fmt.Errorf("error opening bootstrap log file: %s", err)
	}
	fmt.Printf("setting bootstrap log level: '%s'\n", c.String("log-level"))
	if bslog.Level, err = helpers.HigherLogLevel("error", c.String("log-level")); err != nil {
		return fmt.Errorf("error setting bootstrap log level: %s", err)
	}
	bslog.Info("bootstrap log initialized")
	return nil
}

func loadConfig(c *cli.Context) error {
	bslog.Info("loading configuration...")

	var configName = c.String("config")
	var configExtension = c.String("config-ext")
	var configDirectory = c.String("config-dir")

	bslog.WithFields(logrus.Fields{
		"cName": configName,
		"cExt":  configExtension,
		"cDir":  configDirectory,
	}).Info("configuration details")

	viper.SetConfigName(configName)
	viper.SetConfigType(configExtension)
	viper.AddConfigPath(configDirectory)

	bslog.Info("reading configuration")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading in config file: %s", err)
	}
	return nil
}

func loadLog(c *cli.Context) error {
	bslog.Info("loading log")

	var configLog = viper.GetString("log")
	var configLevel = viper.GetString("log-Level")
	var cliLevel = c.String("log-level")

	bslog.WithFields(logrus.Fields{
		"configLog":   configLog,
		"configLevel": configLevel,
		"cliLevel":    cliLevel,
	}).Info("log configuration")

	var err error
	bslog.Info("loading log file")
	if log.Out, err = helpers.AppendOrCreateFile(configLog); err != nil {
		return fmt.Errorf("error while opening log file: %s", err)
	}
	bslog.Info("setting log level")
	if log.Level, err = helpers.HigherLogLevel(configLevel, cliLevel); err != nil {
		return fmt.Errorf("error determing log level: %s", err)
	}
	bslog.Info("log initialized")
	return nil
}

func buildHandlers() {
	bslog.Info("building handlers")
	routeHandlers["/"] = handlers.IndexHandler{
		Handler: handlers.Handler{
			Name:  "IndexHandler",
			Route: "/",
			Log:   log,
		},
	}
}

func buildRouter() *mux.Router {
	bslog.Info("building router")
	r := mux.NewRouter()
	for key, value := range routeHandlers {
		var name = value.GetName()
		bslog.Infof("handling route '%s' with '%s'", key, name)
		r.HandleFunc(key, value.ServeHTTP)
	}
	return r
}

func buildServer(router http.Handler) *http.Server {
	bslog.Info("building server")
	return &http.Server{
		Handler:      router,
		Addr:         viper.GetString("server.address"),
		WriteTimeout: time.Duration(viper.GetInt("server.write-time-out")) * time.Second,
		ReadTimeout:  time.Duration(viper.GetInt("server.read-time-out")) * time.Second,
	}
}
