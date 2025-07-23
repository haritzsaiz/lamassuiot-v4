package main

import (
	"context"
	"log"
	"os"

	"github.com/lamassuiot/lamassuiot/v4/internal/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/controllers"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

	conf, err := config.LoadConfig[ca.CAConfig](nil)
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

	initTracer()

	caService, err := ca.AssembleCAService(conf)
	if err != nil {
		logger.Fatalf("could not assemble User Service: %s", err)
	}

	lHttp := logger.SetupLogger(conf.AppConfig.Server.LogLevel, "API", "HTTP Server")

	httpEngine := http.NewFiberApp(lHttp)
	httpGrp := httpEngine.Group("/")
	ca.NewCAHTTPLayer(&httpGrp, *caService)
	_, err = http.RunHttpServer(lHttp, httpEngine, conf.AppConfig.Server, controllers.APIServiceInfo{
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

func initTracer() func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "MonolithicPKI"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: %s", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}
