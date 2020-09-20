package init

import "ccg-api/docs"

func Swagger() {
	docs.SwaggerInfo.Title = "Swagger CCG-SERVICE API"
	docs.SwaggerInfo.Description = "This is Gola CCG API Server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
}
