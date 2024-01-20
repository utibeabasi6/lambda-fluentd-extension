package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/utibeabasi6/lambda-fluentd-extension/runtime"
)

const (
	INITIAL_QUEUE_SIZE = 5
)

func main() {
	extensionName := path.Base(os.Args[0])
	printPrefix := fmt.Sprintf("[%s]", extensionName)
	logger := log.WithFields(log.Fields{"agent": extensionName})

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-sigs
		cancel()
		logger.Info(printPrefix, "Received", s)
		logger.Info(printPrefix, "Exiting")
	}()

	extensionClient := runtime.ExtensionClient{
		URL:        os.Getenv("AWS_LAMBDA_RUNTIME_API"),
		Version:    "2020-01-01",
		HttpClient: &http.Client{},
		Ctx:        ctx,
	}

	// Register lambda extension
	registrationResponse, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		panic(err)
	}

	logger.Info(printPrefix, "Registered extension for function: "+registrationResponse.FunctionName)

	logChan := make(chan string)
	exitChan := make(chan string, 1)

	// telemetryListener := telemetryapi.TelemetryListener{
	// 	HttpServer: &http.Server{},
	// }

	go func(logChan, exitChan *chan string) {
		for {
			select {
			case <-*exitChan:
				break
			case log := <-*logChan:
				// Send to fluentd
				fmt.Println("recieved log:", log)
			}
		}
	}(&logChan, &exitChan)

	for {
		select {
		case <-ctx.Done():
			exitChan <- "exited"
			return
		default:
			logger.Info(printPrefix, "Waiting for next event")
			// This is a blocking call
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				logger.Info(printPrefix, "Error while fetching next event:", err)
				continue
			}
			// Flush log queue in here after waking up
			// pushLogs(false)
			// Exit if we receive a SHUTDOWN event
			if res.EventType == runtime.ShutDown {
				logger.Info(printPrefix, "Received SHUTDOWN event")
				// pushLogs(true)
				logger.Info(printPrefix, "Exiting")
				return
			}
		}
	}
}
