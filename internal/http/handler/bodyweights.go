package handler

import (
	"github.com/labstack/echo"
)

func BodyWeightCreate(c echo.Context) error {
	return nil
}

func BodyWeightFindOne(c echo.Context) error {
	id := c.Param("id")
	c.String(200, id)
	return nil
}

func BodyWeightFindAll(c echo.Context) error {
	return nil
}

func BodyWeightUpdate(c echo.Context) error {
	return nil
}

func BodyWeightDelete(c echo.Context) error {
	return nil
}
