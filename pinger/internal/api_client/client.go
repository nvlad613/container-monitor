package api_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pinger/config"
	"pinger/internal/container_pinger"
)

const (
	apiRequestPath = "/container/health"
)

type ApiClient interface {
	Post(report container_pinger.HealthReport, ctx context.Context) error
}

type ApiClientImpl struct {
	inner   http.Client
	apiConf config.ServerConnection
}

func New(
	httpClient http.Client,
	apiConf config.ServerConnection,
) *ApiClientImpl {
	return &ApiClientImpl{
		inner:   httpClient,
		apiConf: apiConf,
	}
}

func (c *ApiClientImpl) Post(report container_pinger.HealthReport, _ context.Context) error {
	requestView := new(HealthReport).FromModel(report)
	requestBody, err := json.Marshal(requestView)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	var response *http.Response
	for i := 0; i < c.apiConf.MaxRetries; i++ {
		reader := bytes.NewReader(requestBody)
		response, err = c.inner.Post(
			fmt.Sprintf("%s/%s", c.apiConf.Host, apiRequestPath),
			"application/json",
			reader,
		)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return fmt.Errorf("failed to send report: %w", err)
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send report: response with code %d", response.StatusCode)
	}

	return nil
}
