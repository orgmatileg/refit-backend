package bodyweight

import (
	"github.com/labstack/echo"
)

type IBodyWeight interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

type bodyweight struct{}

// New Repository todos
func New() IBodyWeight {
	return &bodyweight{}
}

func (u bodyweight) Create(c echo.Context) error {
	return nil
}
func (u bodyweight) FindOne(c echo.Context) error {
	return nil
}
func (u bodyweight) FindAll(c echo.Context) error {
	return nil
}
func (u bodyweight) Update(c echo.Context) error {
	return nil
}
func (u bodyweight) Delete(c echo.Context) error {
	return nil
}
func (u bodyweight) Count(c echo.Context) error {
	return nil
}
