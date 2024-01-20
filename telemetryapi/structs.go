package telemetryapi

type TelemetryEvent string

type URI string

type Destination struct {
	Protocol   HttpProtocol `json:"protocol"`
	URI        URI          `json:"URI"`
	HttpMethod HttpMethod   `json:"method"`
	Encoding   HttpEncoding `json:"encoding"`
}

type SchemaVersion string

const (
	SchemaVersion20220701 = "2022-07-01"
	SchemaVersionLatest   = SchemaVersion20220701
)

// Request body that is sent to the Telemetry API on subscribe
type SubscribeRequest struct {
	SchemaVersion SchemaVersion    `json:"schemaVersion"`
	EventTypes    []TelemetryEvent `json:"types"`
	BufferingCfg  BufferingCfg     `json:"buffering"`
	Destination   Destination      `json:"destination"`
}

// Response body that is received from the Telemetry API on subscribe
type SubscribeResponse struct {
	body string
}

type BufferingCfg struct {
	// Maximum number of log events to be buffered in memory. (default: 10000, minimum: 1000, maximum: 10000)
	MaxItems uint32 `json:"maxItems"`
	// Maximum size in bytes of the log events to be buffered in memory. (default: 262144, minimum: 262144, maximum: 1048576)
	MaxBytes uint32 `json:"maxBytes"`
	// Maximum time (in milliseconds) for a batch to be buffered. (default: 1000, minimum: 100, maximum: 30000)
	TimeoutMS uint32 `json:"timeoutMs"`
}

type HttpProtocol string

type HttpEncoding string

type HttpMethod string

const (
	// Used to receive log events emitted by the platform
	Platform TelemetryEvent = "platform"
	// Used to receive log events emitted by the function
	Function TelemetryEvent = "function"
	// Used is to receive log events emitted by the extension
	Extension TelemetryEvent = "extension"

	JSON HttpEncoding = "JSON"

	HttpPost HttpMethod = "POST"
	// Receive log events via PUT requests to the listener
	HttpPut HttpMethod = "PUT"

	HttpProto HttpProtocol = "HTTP"
)
