package main

import (
	"flag"
	"log"

	"github.com/perun-network/erdstall/operator"
)

func main() {
	configFilePath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	cfg, err := operator.LoadConfig(*configFilePath)
	operator.AssertNoError(err)
	log.Println("Config loaded")

	_operator, _ := operator.Setup(cfg)
	err = _operator.Serve(cfg.Port)
	operator.AssertNoError(err)
}
