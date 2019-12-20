package handler

import (
	"github.com/labstack/echo"
)

func TodosCreate(c echo.Context) error {
	return nil
}

func TodosFindOne(c echo.Context) error {
	id := c.Param("id")
	c.String(200, id)
	return nil
}

func TodosFindAll(c echo.Context) error {
	return nil
}

func TodosUpdate(c echo.Context) error {
	return nil
}

func TodosDelete(c echo.Context) error {
	return nil
}
