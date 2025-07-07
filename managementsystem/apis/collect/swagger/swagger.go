package swagger

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func New(e *echo.Echo) {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "management system hr APIs"
	docs.SwaggerInfo.Description = "Management system for HR Provider REST APIs"

	e.GET("/v1/swagger/*", echoSwagger.WrapHandler)
}
