package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/cryptoutils"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

func BuildHTTPClientWithTLSOptions(cli *http.Client, cfg config.TLSConfig) (*http.Client, error) {
	caPool := cryptoutils.LoadSytemCACertPool()
	tlsConfig := &tls.Config{}

	if cfg.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	if cfg.CACertificateFile != "" {
		cert, err := cryptoutils.ReadCertificateFromFile(cfg.CACertificateFile)
		if err != nil {
			return nil, err
		}

		caPool.AddCert(cert)
	}

	cli.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return cli, nil
}

func BuildHTTPClientWithTracerLogger(cli *http.Client, logger *logger.Logger) (*http.Client, error) {
	transport := http.DefaultTransport
	if cli.Transport != nil {
		transport = cli.Transport
	}

	cli.Transport = loggingRoundTripper{
		transport: transport,
		logger:    logger,
	}

	return cli, nil
}

type loggingRoundTripper struct {
	transport http.RoundTripper
	logger    *logger.Logger
}

func (lrt loggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	start := time.Now()
	// Send the request, get the response (or the error)
	res, err = lrt.transport.RoundTrip(req)
	if err != nil {
		lrt.logger.Errorf("%s: %s", req.URL.String(), err)
	} else {
		log := lrt.logger.With("response", fmt.Sprintf("%s %d: %s", req.Method, res.StatusCode, time.Since(start)))
		log.Debug(req.URL.String())
		dumpReq, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Error("Failed to dump request:", err)
			dumpReq = []byte("<< Failed to dump request >>")
		}

		dumpResp, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Error("Failed to dump response:", err)
			dumpResp = []byte("<< Failed to dump response >>")
		}

		log.Tracef("--- Request ---\n%s\n--- Response ---\n%s\n", string(dumpReq), string(dumpResp))
	}

	return
}
