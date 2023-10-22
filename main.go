package main

import (
	. "ccg-api/init"
	"context"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/tracing"
	"net/http"
	"strings"
)

func main() {
	logger := logging.NewLoggerEntry()
	logger.Info("Starting service...")

	logger.Info("Loading configurations")
	configData := LoadConfig()
	router := CreateRouter(configData)
	_, _ = tracing.Init(configData.TracingServiceName, configData.TracingOCAgentHost)
	var port string
	if strings.EqualFold(configData.Environment, "local") {
		port = ":8083"
	} else {
		port = ":8080"
	}
	err := http.ListenAndServe(port, tracing.WithTracing(router, "/api/ccg/healthz"))
	if err != nil {
		logging.GetLogger(context.TODO()).Error("Could not start the server", err)
	}
}
