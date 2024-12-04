package main

import (
	"flag"
	"os"
	"time"
	"wallet/internal/app"
	"wallet/internal/config"
	"wallet/pkg/logger"
)

func main() {
	time.Sleep(10 * time.Second)
	log := logger.SetupLogger(logger.Local, "./logs.log")
	configPath := flag.String("c", "", "Path to the configuration file")
	flag.Parse()

	if *configPath != "" {
		log.Info("Trying to load configuration from", "file", *configPath)
	}
	cfg := config.MustLoad(*configPath)
	a, err := app.NewApp(log, cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	a.Start()
}
