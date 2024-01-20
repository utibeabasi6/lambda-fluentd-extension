package runtime

import (
	"context"
	"net/http"
)

type ExtensionClient struct {
	URL         string
	Version     string
	HttpClient  *http.Client
	Ctx         context.Context
	ExtensionId string `json:"extensionId"`
}

type RegistrationResponse struct {
	FunctionName    string `json:"functionName"`
	FunctionVersion string `json:"functionVersion"`
	Handler         string `json:"handler"`
}

type RegistrationEvent string

const (
	Invoke                   RegistrationEvent = "INVOKE"
	ShutDown                 RegistrationEvent = "SHUTDOWN"
	extensionNameHeader                        = "Lambda-Extension-Name"
	extensionIdentiferHeader                   = "Lambda-Extension-Identifier"
)

type NextEventResponse struct {
	EventType          RegistrationEvent `json:"eventType"`
	DeadlineMs         int64             `json:"deadlineMs"`
	RequestID          string            `json:"requestId"`
	InvokedFunctionArn string            `json:"invokedFunctionArn"`
	Tracing            Tracing           `json:"tracing"`
}

type Tracing struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
