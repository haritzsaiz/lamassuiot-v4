package server

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/controllers"
	fiber_context_mw "github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/middleware/context"
	fiber_logger_mw "github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/middleware/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

func NewFiberApp(log *logger.Logger) *fiber.App {
	fiber_logger_mw.ForceConsoleColor()
	router := fiber.New()

	router.Use(recover.New())

	corsCfg := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}

	TIMEOUT := time.Second * 10

	router.Use(
		fiber_context_mw.WithContext(TIMEOUT),
		cors.New(corsCfg),
		otelfiber.Middleware(),
		fiber_logger_mw.UseLogger(log),
		fiber_logger_mw.DumpWithOptions(true, true, true, true, func(dumpStr string) {
			log.Trace(dumpStr)
		}),
	)

	return router
}

func RunHttpServer(log *logger.Logger, routerEngine *fiber.App, httpServerCfg config.HttpServer, apiInfo controllers.APIServiceInfo) (int, error) {
	hCheckRoute := controllers.NewHealthCheckRoute(apiInfo)
	mainLogger := log
	if !httpServerCfg.HealthCheckLogging {
		mainLogger = logger.New(logger.NewFormatterHandler(io.Discard, logger.LevelInfo, logger.LogFormatter))
	}

	healthEngine := NewFiberApp(mainLogger)
	healthEngine.Get("/health", hCheckRoute.HealthCheck)

	mainEngine := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	mainEngine.Mount("/", routerEngine)
	mainEngine.Mount("/health", healthEngine)

	addr := fmt.Sprintf("%s:%d", httpServerCfg.ListenAddress, httpServerCfg.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", httpServerCfg.ListenAddress, httpServerCfg.Port))
	if err != nil {
		return -1, err
	}

	usedPort := listener.Addr().(*net.TCPAddr).Port

	wg := new(sync.WaitGroup)
	wg.Add(1) // add `1` goroutines to finish
	startLaunching := func() {
		wg.Done()
	}

	for _, routes := range routerEngine.Stack() {
		for _, route := range routes {
			if len(route.Handlers) > 1 {
				continue
			}
			logger.Debug(fmt.Sprintf("Endpoint: %-6s %s", route.Method, route.Path))
		}
	}

	httpErrChan := make(chan error, 1)

	go func() {
		if httpServerCfg.Protocol == config.HTTPS {
			logger.Info(fmt.Sprintf("HTTPS server listening on %s:%d", addr, usedPort))
			startLaunching()

			cert, err := tls.LoadX509KeyPair(httpServerCfg.CertFile, httpServerCfg.KeyFile)
			if err != nil {
				log.Fatalf("failed to load TLS cert: %v", err)
				httpErrChan <- err
			}

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
			tlsListener := tls.NewListener(listener, tlsConfig)

			err = mainEngine.Listener(tlsListener)

			if err == nil {
				logger.Error(fmt.Sprintf("could not start http server: %s", err))
				httpErrChan <- err
			}
		} else {
			logger.Info(fmt.Sprintf("HTTP server listening on %s", addr))
			startLaunching()

			err := mainEngine.Listener(listener)
			if err == nil {
				logger.Error(fmt.Sprintf("could not start http server: %s", err))
				httpErrChan <- err
			}
		}
	}()

	return usedPort, nil
}
