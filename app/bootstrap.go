package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jeff-All/Dohyo/authentication"
	"github.com/Jeff-All/Dohyo/handlers"
	"github.com/Jeff-All/Dohyo/helpers"
	"github.com/Jeff-All/Dohyo/middlewares"
	"github.com/Jeff-All/Dohyo/services"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var bslog = logrus.New()
var routeHandlers = make(map[string]handlers.HandlerInterface)
var middleware = make(map[string]middlewares.MiddlewareInterface)
var rankService services.RankService
var rikishiService services.RikishiService
var categoryService services.CategoryService
var teamService services.TeamService
var matchService services.MatchService
var tournamentService services.TournamentService
var db *gorm.DB
var log = logrus.New()

func bootstrap(c *cli.Context) error {
	if err := loadBootstrapLog(c); err != nil {
		panic(fmt.Errorf("error loading bootstrap log: %s", err))
	}
	if err := loadConfig(c); err != nil {
		bslog.Errorf("error loading config: %s", err)
	}
	if err := loadLog(c); err != nil {
		bslog.Errorf("error loading log: %s", err)
	}
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
	if c.Bool("log-clear") {
		bslog.Infof("deleting log file")
		if err = helpers.DeleteFile(configLog); err != nil {
			return fmt.Errorf("error while deleting log file: %s", err)
		}
	}
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

func loadDB() error {
	bslog.Info("connecting to database")
	filename := viper.GetString("db.filename")
	if err := helpers.CreateDirectoryForFile(filename); err != nil {
		return err
	}
	var err error
	if db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{}); err != nil {
		return err
	}
	bslog.Info("connected to database")
	return nil
}

func buildAuthentication() error {
	bslog.Info("building authentication")
	return authentication.SetNewAuthenticator(log, viper.GetString("authentication.config"))
}

func buildServices() {
	tournamentService = services.NewTournamentService(log, db)
	rankService = services.NewRankService(log, db)
	rikishiService = services.NewRikishiService(log, db, rankService, tournamentService)
	categoryService = services.NewCategoryService(log, db, rikishiService)
	teamService = services.NewTeamService(log, db)
	matchService = services.NewMatchService(log, db, rikishiService, tournamentService)
}

func buildHandlers() {
	bslog.Info("building handlers")

	routeHandlers["index"] = handlers.IndexHandler{
		Handler: handlers.Handler{
			Log: log,
		},
	}

	routeHandlers["rikishis"] = handlers.RikishisHandler{
		Handler: handlers.Handler{
			Log: log,
		},
		RikishiService: rikishiService,
	}

	routeHandlers["categorizedRikishis"] = handlers.CategorizedRikishiHandler{
		Handler: handlers.Handler{
			Log: log,
		},
		CategoryService: categoryService,
	}

	routeHandlers["teams"] = handlers.TeamHandler{
		Handler: handlers.Handler{
			Log: log,
		},
		TeamService:     teamService,
		CategoryService: categoryService,
	}

	bslog.Info("handlers built")
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

func buildRouter() *mux.Router {
	bslog.Info("building router")
	r := mux.NewRouter()

	return r
}

func buildMiddleware() {
	authorizationCacheDuration := viper.GetDuration("authorization.cacheDuration")
	authorizationCacheCleanup := viper.GetDuration("authorization.cacheCleanup")
	bslog.Infof(
		"building authorization middleware(cacheDuration=%v,cacheCleanup=%v",
		authorizationCacheDuration,
		authorizationCacheCleanup,
	)
	middleware["authorization"] = &middlewares.AuthorizationMiddleware{
		Log: log,
		Cache: cache.New(
			authorizationCacheDuration,
			authorizationCacheCleanup,
		),
	}
}
