package main

import (
	"OTUS_hws/Anti-BruteForce/internal/antibrutforce"
	"OTUS_hws/Anti-BruteForce/internal/config"
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi"
	"OTUS_hws/Anti-BruteForce/internal/gen/restapi/operations"
	"OTUS_hws/Anti-BruteForce/internal/handlers"
	"OTUS_hws/Anti-BruteForce/internal/redisdb"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/dev.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "[config.New()]")
		panic(err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		err = errors.Wrap(err, "[loads.Analyzed()]")
		panic(err)
	}

	redisClient := redisdb.NewClient(*conf)
	abfChecker := antibrutforce.New(redisClient, conf)
	h := handlers.NewHandler(abfChecker)

	api := operations.NewAntiBrutForceAPI(swaggerSpec)
	h.Register(api)
	server := restapi.NewServer(api)
	server.Port = conf.Service.Port
	server.Host = conf.Service.Host

	if err = server.Serve(); err != nil {
		err = errors.Wrap(err, "[server.Serve()]")
		panic(err)
	}

	sig := <-quit

	err = redisClient.Client.Close()
	if err != nil {
		err = errors.Wrapf(err, "[db.Close(%v)]", sig)
		panic(err)
	}

	err = server.Shutdown()
	if err != nil {
		err = errors.Wrap(err, "[server.Shutdown()]")
		panic(err)
	}
}
