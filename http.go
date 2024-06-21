package slogger

import (
	"encoding/json"
	"net/http"
)

// TODO 2024-06-21 Sam Borkent: add tests

var _ http.Handler = &Logger{}

// Request to change the program log level at runtime through a HTTP endpoint.
type LogLevelRequest struct {
	LogLevel string `json:"logLevel"`
}

// Handler function to enable changing the program log level at runtime through a HTTP endpoint.
func (l *Logger) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// The progam level gets set, so we should only accept PUT calls.
	if req.Method != http.MethodPut {
		http.Error(writer, "method not supported, expected PUT", http.StatusBadRequest)
		return
	}

	// Read the request body.
	bodyData := make([]byte, 512)

	nBytesRead, err := req.Body.Read(bodyData)
	if err != nil {
		http.Error(writer, "reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if nBytesRead == 0 {
		http.Error(writer, "empty request", http.StatusBadRequest)
		return
	}

	// Parse the log level request.
	var logLevelRequest LogLevelRequest

	if err := json.Unmarshal(bodyData[:nBytesRead], &logLevelRequest); err != nil {
		http.Error(writer, "parsing request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the log level string to a slog.Level.
	logLevel, err := stringToSlogLevel(logLevelRequest.LogLevel)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the global program log level.
	programLevel.Set(logLevel)
}
