package main

import (
	. "ccg-api/init"
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/tracing"
	"net/http"
)

func main() {
	logger := logging.NewLoggerEntry()
	logger.Info("Starting service...")

	logger.Info("Loading configurations")
	configData := LoadConfig()
	router := CreateRouter(configData)
	tracing.Init(configData.TracingServiceName, configData.TracingOCAgentHost)
	var port string
	if configData.Environment == "local" {
		port = ":8080"
	} else {
		port = ":8083"
	}
	err := http.ListenAndServe(port, tracing.WithTracing(router, "/api/ccg/healthz"))
	if err != nil {
		logging.GetLogger(context.TODO()).Error("Could not start the server", err)
	}
}
