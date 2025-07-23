package kms

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type KMSSdkService struct{}

func NewKMSSdkService() *KMSSdkService {
	return &KMSSdkService{}
}

func (s *KMSSdkService) CreateKMSKey(ctx context.Context, input CreateKMSInput) error {
	ctx, span := otel.GetTracerProvider().Tracer("kms-sdk").Start(ctx, "CreateKMSKey", trace.WithAttributes(semconv.PeerService("KMS")))
	defer span.End()

	body := map[string]string{
		"name": input.Name,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		span.RecordError(err)
		return err
	}
	byteReader := bytes.NewReader(jsonBody)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8091/v1/kms", byteReader)
	if err != nil {
		span.RecordError(err)
		return err
	}

	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
		),
	}

	res, err := client.Do(r)
	if err != nil {
		span.RecordError(err)
		return err
	}

	defer res.Body.Close()

	// Record the HTTP status code
	span.SetAttributes(semconv.HTTPStatusCode(res.StatusCode))

	if res.StatusCode >= 400 {
		span.SetAttributes(semconv.HTTPStatusCode(res.StatusCode))
		// You might want to read the response body for error details
	}

	return nil
}

func (s *KMSSdkService) GetKMSKeys(ctx context.Context, input GetKMSKeysInput) (string, error) {
	// Implementation for retrieving KMS keys
	return "", nil
}
