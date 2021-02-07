package app

import (
	"github.com/Jeff-All/Dohyo/helpers"
	"github.com/Jeff-All/Dohyo/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := buildAuthentication(); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}
	buildServices()
	buildHandlers()
	router := buildRouter()
	defineRoutes(router)
	server := buildServer(router)
	bslog.Info("launching server")
	log.Fatal(server.ListenAndServe())
	return nil
}

func load(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}
	log.Info("loading data into database")
	dataFile := viper.New()
	dir, name, ext := helpers.SplitFileName(c.String("data-file"))

	log.WithFields(logrus.Fields{
		"dir":  dir,
		"name": name,
		"ext":  ext,
	}).Info("data file details")

	dataFile.SetConfigName(name)
	dataFile.AddConfigPath(dir)
	dataFile.SetConfigType(ext)

	if err := dataFile.ReadInConfig(); err != nil {
		return err
	}

	buildServices()
	service := services.NewLoadService(log, db, dataFile, rankService, rikishiService, categoryService)

	for _, arg := range c.Args().Slice() {
		if err := service.Load(arg); err != nil {
			log.Errorf("error while loading %s: %s", arg, err)
			return err
		}
	}
	return nil
}

func migrate(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}
	log.Info("migrating tables into database")

	service := services.NewMigrationService(log, db)
	service.MigrateModels(c.Args().Slice()...)
	return nil
}
