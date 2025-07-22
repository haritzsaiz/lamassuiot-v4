package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/middleware/context"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

const (
	defaultLogFormat = "%s %-7s %s %s %3d %s [ %30v ] | %13v | \"%s\"\n"
)

type traceRequestWriter struct {
	logger *logger.Logger
}

func (tr *traceRequestWriter) Write(p []byte) (n int, err error) {
	tr.logger.Debug(fmt.Sprintf("%s", string(p)))
	return len(p), nil
}

func UseLogger(log *logger.Logger) fiber.Handler {
	return logRequest(log)
}

func loggingWithReqBodyLog(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	return fmt.Sprintf(defaultLogFormat,
		methodColor, param.Method, resetColor,
		statusColor, param.StatusCode, resetColor,
		fmt.Sprintf("agent: \"%s\"", string(param.Context.Request().Header.UserAgent())),
		param.Latency,
		param.Path,
	)
}

func logRequest(l *logger.Logger) fiber.Handler {
	formatter := loggingWithReqBodyLog

	return func(c *fiber.Ctx) error {
		request := c.Request()
		resp := c.Response()
		ctx := context.GetRequestContext(c)
		lReq := logger.ConfigureLogger(ctx, l)
		out := &traceRequestWriter{logger: lReq}

		// Start timer
		start := time.Now()
		path := c.Path()
		raw := string(request.URI().QueryString())

		// Process request
		err := c.Next()

		param := LogFormatterParams{
			Context: c,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.IP()
		param.Method = string(request.Header.Method())
		param.StatusCode = resp.StatusCode()
		//param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = len(resp.Body())

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		fmt.Fprint(out, formatter(param))
		return err
	}
}

// LogFormatter
type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

var consoleColorMode = autoColor

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type LogFormatterParams struct {
	Context *fiber.Ctx

	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
	isTerm       bool
	BodySize     int
	Keys         map[string]any
}

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func (p *LogFormatterParams) ResetColor() string {
	return reset
}

func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

func ForceConsoleColor() {
	consoleColorMode = forceColor
}
