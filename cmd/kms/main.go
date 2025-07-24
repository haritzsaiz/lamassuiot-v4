package main

import (
	"os"

	"github.com/lamassuiot/lamassuiot/v4/internal/kms"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/controllers"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"gopkg.in/yaml.v2"
)

var (
	version   string = "v0"    // api version
	sha1ver   string = "-"     // sha1 revision used to build the program
	buildTime string = "devTS" // when the executable was built
)

func main() {
	var logLevel = new(logger.LevelVar)

	log := logger.New(logger.NewFormatterHandler(os.Stdout, logLevel, logger.LogFormatter))

	logger.SetDefault(log)
	logger.Infof("starting api: version=%s buildTime=%s sha1ver=%s", version, buildTime, sha1ver)

	conf, err := config.LoadConfig[kms.KMSConfig](nil)
	if err != nil {
		logger.Fatalf("something went wrong while loading config. Exiting: %s", err)
	}

	logLevel.Set(conf.AppConfig.Logs.Level)
	logger.Infof("global log level set to '%s'", conf.AppConfig.Logs.Level)

	confBytes, err := yaml.Marshal(conf)
	if err != nil {
		logger.Fatalf("could not dump yaml config: %s", err)
	}

	logger.Debug("===================================================")
	logger.Debug(string(confBytes))
	logger.Debug("===================================================")

	kmsService, err := kms.AssembleKMSService(conf)
	if err != nil {
		logger.Fatalf("could not assemble KMS Service: %s", err)
	}

	lHttp := logger.SetupLogger(conf.AppConfig.Server.LogLevel, "API", "HTTP Server")

	httpEngine := server.NewFiberApp(lHttp)
	httpGrp := httpEngine.Group("/")

	kms.NewKMSHTTPLayer(&httpGrp, *kmsService)

	_, err = server.RunHttpServer(lHttp, httpEngine, conf.AppConfig.Server, controllers.APIServiceInfo{
		Version:   version,
		BuildSHA:  sha1ver,
		BuildTime: buildTime,
	})

	if err != nil {
		logger.Fatalf("could not run API Server. Exiting: %s", err)
	}

	forever := make(chan struct{})
	<-forever
}
