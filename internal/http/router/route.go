package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"refit_backend/internal/http/handler"
)

// SetHTTPRouter ...
func SetHTTPRouter(e *echo.Echo) {

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/healthcheck", handler.HealthCheck)

	routerUsers := e.Group("users")
	routerUsers.GET("", handler.HealthCheck)
}
