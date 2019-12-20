package handler

import (
	"github.com/labstack/echo"
)

func CountriesFindOne(c echo.Context) error {
	id := c.Param("id")
	c.String(200, id)
	return nil
}

func CountriesFindAll(c echo.Context) error {
	return nil
}
