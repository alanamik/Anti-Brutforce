package main

import (
	"OTUS_hws/Anti-BruteForce/internal/antibrutforce"
	"OTUS_hws/Anti-BruteForce/internal/config"
	server "OTUS_hws/Anti-BruteForce/internal/server/http"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "[config.New()]")
		panic(err)
	}

	abfChecker, err := antibrutforce.New(conf)
	if err != nil {
		panic(err)
	}
	server := server.New(abfChecker, conf)

	interval := time.Duration(10) * time.Minute // clear one time in 10 minutes
	tk := time.NewTicker(interval)
	tickerChan := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerChan:
				return
			case <-tk.C:
				abfChecker.ClearOldBuckets()
			}
		}
	}()

	go func() {
		if err := server.Start(); err != nil {
			os.Exit(1)
		}
	}()
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Stop(ctx)
}
