package handler

import (
	"github.com/labstack/echo"
)

func UsersCreate(c echo.Context) error {
	return nil
}

func UsersFindOne(c echo.Context) error {
	id := c.Param("id")
	c.String(200, id)
	return nil
}

func UsersFindAll(c echo.Context) error {
	return nil
}

func UsersUpdate(c echo.Context) error {
	return nil
}

func UsersDelete(c echo.Context) error {
	return nil
}
