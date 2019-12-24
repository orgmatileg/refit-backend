package tools

import (
	"net/http"
	"refit_backend/internal/helpers"
	"strings"

	"github.com/labstack/echo"
)

type ITools interface {
	HealthCheck(c echo.Context) error
	DefaultErrorHandler(err error, c echo.Context)
}

type tools struct{}

// New Repository Users
func New() ITools {
	return &tools{}
}

func (u tools) HealthCheck(c echo.Context) error {
	return nil
}

func (u tools) DefaultErrorHandler(err error, c echo.Context) {
	var (
		code int = http.StatusInternalServerError
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if err != nil {
		if strings.Contains(err.Error(), "not exists") {
			code = 404
		}
		helpers.MakeDefaultResponse(c, code, err)
	}
}
