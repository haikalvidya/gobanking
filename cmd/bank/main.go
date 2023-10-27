package main

import (
	"flag"
	"fmt"
	"gobanking/internal/bank/config"
	"gobanking/pkg/logger"
	"log"
	"strings"
)

func main() {
	log.Println("Starting the application...")

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("config.InitConfig: %v", err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Named(fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName)))
	appLogger.Infof("CFG: %+v", cfg)
}
