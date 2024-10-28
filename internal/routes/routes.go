package routes

import (
	"log/slog"
	"net/http"

	"github.com/jha-captech/blog/internal/handlers"
	"github.com/jha-captech/blog/internal/services"
	"github.com/swaggo/http-swagger/v2"

	_ "github.com/jha-captech/blog/cmd/api/docs"
)

// AddRoutes adds all routes to the provided mux.
//
//	@title						Blog Service API
//	@version					1.0
//	@description				Practice Go Gin API using GORM and Postgres
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8000
//	@BasePath					/api
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, baseURL string) {
	// Read a user
	mux.Handle("GET /api/users/{id}", handlers.HandleReadUser(logger, usersService))

	// swagger docs
	mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))
}
