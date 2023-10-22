package init

import (
	"ccg-api/configuration"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/middleware/request_response_trace"
	middleware "github.com/inclusi-blog/gola-utils/middleware/session_trace"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRouter(router *gin.Engine, configData *configuration.ConfigData) {
	routerGroup := router.Group("/api")
	routerGroup.GET("/ccg/healthz", healthController.GetHealth)
	router.Use(middleware.SessionTracingMiddleware)
	router.Use(request_response_trace.HttpRequestResponseTracingMiddleware([]request_response_trace.IgnoreRequestResponseLogs{
		{
			PartialApiPath:       "api/ccg/healthz",
			IsRequestLogAllowed:  false,
			IsResponseLogAllowed: false,
		},
	}, "api/ccg/healthz", nil, nil))

	golaLoggerRegistry := logging.NewLoggerEntry()

	router.Use(logging.LoggingMiddleware(golaLoggerRegistry))

	logLevel := configData.LogLevel
	logger := logging.GetLogger(context.TODO())

	if logLevel != "" {
		logLevelInitErr := golaLoggerRegistry.SetLevel(logLevel)
		if logLevelInitErr != nil {
			logger.Warning("gola_logger.SetLevel failed. Default log level being used", logLevelInitErr.Error())
		}
	}

	router.GET("api/ccg/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		routerGroup.POST("/ccg/v1/email/send", emailController.SendEmail)
	}

}
