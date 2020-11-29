package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/operator"
)

func main() {
	configFilePath := flag.String("config", "config.json", "config file path")
	logLevel := flag.String("log-level", "info", "log level")
	flag.Parse()

	cfg, err := operator.LoadConfig(*configFilePath)
	operator.AssertNoError(err)
	log.Println("Config loaded")

	lvl, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("parsing log level: %v", err)
	}
	log.SetLevel(lvl)

	_operator := operator.Setup(*cfg)
	err = _operator.Serve(cfg.Port)
	operator.AssertNoError(err)
}
