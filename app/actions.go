package app

import (
	"github.com/Jeff-All/Dohyo/helpers"
	"github.com/Jeff-All/Dohyo/scrapers"
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
	buildMiddleware()
	buildHandlers()
	router := buildRouter()
	defineRoutes(router)
	server := buildServer(router)
	bslog.Info("launching server")
	log.Fatal(server.ListenAndServe())
	return nil
}

func activate(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}
	buildServices()
	return tournamentService.SetCurrentTournament(c.String("name"))
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
	service := services.NewLoadService(log, db, dataFile, rankService, rikishiService, categoryService, matchService)

	for _, arg := range c.Args().Slice() {
		if err := service.Load(arg, c.Bool("clear")); err != nil {
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

func sql(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}
	var err error
	var sql []byte
	filename := c.String("file")
	if sql, err = helpers.ReadFile(filename); err != nil {
		log.Errorf("error reading file %s", filename)
	}
	log.Info("executing SQL against the database: %s", string(sql))
	if err = db.Exec(string(sql)).Error; err != nil {
		log.Errorf("error executing SQL: %s")
		return err
	}
	return nil
}

func scrapeMatches(c *cli.Context) error {
	if err := bootstrap(c); err != nil {
		return err
	}
	if err := loadDB(); err != nil {
		return err
	}

	scraper := scrapers.NewMatchScraper(log, viper.GetString("scrapers.match.url"), viper.GetString("scrapers.match.host"))
	service := services.NewScraperService(log, db, scraper)

	return service.ScrapeMatches(c.String("month"), c.Uint("day"))
}
