package handler

import (
	"github.com/labstack/echo"
	"net/http"
)

// HealthCheck handler
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
