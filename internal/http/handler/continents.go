package handler

import (
	"github.com/labstack/echo"
)

func ContinentFindOne(c echo.Context) error {
	id := c.Param("id")
	c.String(200, id)
	return nil
}

func ContinentFindAll(c echo.Context) error {
	return nil
}
