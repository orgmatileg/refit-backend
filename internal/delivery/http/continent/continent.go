package continent

import (
	"github.com/labstack/echo"
)

type IContinent interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

type continent struct{}

// New Repository Users
func New() IContinent {
	return &continent{}
}

func (co continent) Create(c echo.Context) error {
	return nil
}
func (co continent) FindOne(c echo.Context) error {
	return nil
}
func (co continent) FindAll(c echo.Context) error {
	return nil
}
func (co continent) Update(c echo.Context) error {
	return nil
}
func (co continent) Delete(c echo.Context) error {
	return nil
}
func (co continent) Count(c echo.Context) error {
	return nil
}
