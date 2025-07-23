package ca

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type CASdkService struct{}

func NewCASdkService() *CASdkService {
	return &CASdkService{}
}

func (s *CASdkService) CreateCA(ctx context.Context, input CreateCAInput) error {
	body := map[string]string{
		"name": input.Name,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(jsonBody)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8090/v1/ca", byteReader)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

func (s *CASdkService) GetCAs(ctx context.Context, input GetCAsInput) (string, error) {
	// Implementation for retrieving CAs
	return "", nil
}
