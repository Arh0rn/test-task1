package swagger

import (
	_ "github.com/Arh0rn/test-task1/docs" //To know what to swag.
	"github.com/Arh0rn/test-task1/pkg/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

const swaggerUrl = "/swagger/doc.json"

func Set(cfg *config.HTTPServer) http.HandlerFunc {
	//url := "http://" + cfg.Address + "/swagger/doc.json"
	//return httpSwagger.Handler(httpSwagger.URL(url))
	// No need if we use import docs.
	return httpSwagger.Handler(httpSwagger.URL(swaggerUrl))
}
