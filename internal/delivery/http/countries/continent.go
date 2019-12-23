package countries

import (
	"github.com/labstack/echo"
)

// ICountries interface
type ICountries interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

type countries struct{}

// New con
func New() ICountries {
	return &countries{}
}

func (co countries) Create(c echo.Context) error {
	return nil
}
func (co countries) FindOne(c echo.Context) error {
	return nil
}
func (co countries) FindAll(c echo.Context) error {
	return nil
}
func (co countries) Update(c echo.Context) error {
	return nil
}
func (co countries) Delete(c echo.Context) error {
	return nil
}
func (co countries) Count(c echo.Context) error {
	return nil
}
