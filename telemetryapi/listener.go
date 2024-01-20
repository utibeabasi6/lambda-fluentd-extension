package telemetryapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type TelemetryListener struct {
	HttpServer *http.Server
	LogChannel *chan interface{}
}

func listenOnAddress() string {
	env_aws_local, ok := os.LookupEnv("AWS_SAM_LOCAL")
	var addr string
	if ok && env_aws_local == "true" {
		addr = ":" + telemetryPort
	} else {
		addr = "sandbox:" + telemetryPort
	}

	return addr
}

func (c *TelemetryListener) Start() {
	address := listenOnAddress()
	logger.Info("[TelemetryListener:Start] Starting on address", address)
	http.HandleFunc("/job", c.Handler)
	log.Println("Starting listener on port", telemetryPort)

	http.ListenAndServe(fmt.Sprintf(":%s", telemetryPort), nil)
}

func (c *TelemetryListener) Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("[TelemetryListener:http_handler] Error reading body:", err)
		return
	}

	// Parse and put the log messages into the queue
	var slice []interface{}
	_ = json.Unmarshal(body, &slice)

	counter := 0

	for _, el := range slice {
		*c.LogChannel <- el
		counter += 1
	}

	logger.Info("[listener:http_handler] logEvents received:", counter)
}
