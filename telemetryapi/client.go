package telemetryapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var logger = log.WithFields(log.Fields{"pkg": "telemetryApi"})

type TelemetryClient struct {
	URL         string
	Version     string
	HttpClient  *http.Client
	ExtensionId string `json:"extensionId"`
}

func (c *TelemetryClient) Subscribe(ctx context.Context, extensionId string, listenerUri string) (*SubscribeResponse, error) {
	eventTypes := []TelemetryEvent{
		// Platform,
		Function,
		// Extension,
	}

	bufferingConfig := BufferingCfg{
		MaxItems:  1000,
		MaxBytes:  256 * 1024,
		TimeoutMS: 1000,
	}

	destination := Destination{
		Protocol:   HttpProto,
		HttpMethod: HttpPost,
		Encoding:   JSON,
		URI:        URI(listenerUri),
	}

	data, err := json.Marshal(
		&SubscribeRequest{
			SchemaVersion: SchemaVersionLatest,
			EventTypes:    eventTypes,
			BufferingCfg:  bufferingConfig,
			Destination:   destination,
		})

	if err != nil {
		return nil, errors.New(err.Error() + ": Failed to marshal SubscribeRequest")
	}

	logger.Info("[TelemetryApiClient:Subscribe] Subscribing using baseUrl:", c.URL)

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", c.URL, bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set(lambdaAgentIdentifierHeaderKey, c.ExtensionId)
	httpReq.Header.Set("Content-Type", "application/json")

	httpRes, err := c.HttpClient.Do(httpReq)

	if err != nil {
		logger.Error("[TelemetryApiClient:Subscribe] Subscription failed:", err)
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusAccepted {
		logger.Error("[TelemetryApiClient:Subscribe] Subscription failed. Logs API is not supported! Is this extension running in a local sandbox?")
	} else if httpRes.StatusCode != http.StatusOK {
		logger.Error("[TelemetryApiClient:Subscribe] Subscription failed")
		body, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, errors.Errorf("%s failed: %d[%s]", c.URL, httpRes.StatusCode, httpRes.Status)
		}

		return nil, errors.Errorf("%s failed: %d[%s] %s", c.URL, httpRes.StatusCode, httpRes.Status, string(body))
	}

	body, _ := io.ReadAll(httpRes.Body)
	logger.Info("[TelemetryApiClient:Subscribe] Subscription success:", string(body))

	return &SubscribeResponse{string(body)}, nil

}

func (c *TelemetryClient) Handler(w http.ResponseWriter, r *http.Request) {}
