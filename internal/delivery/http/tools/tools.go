package tools

import (
	"github.com/labstack/echo"
)

type ITools interface {
	HealthCheck(c echo.Context) error
}

type tools struct{}

// New Repository Users
func New() ITools {
	return &tools{}
}

func (u tools) HealthCheck(c echo.Context) error {
	return nil
}
