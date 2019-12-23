package todos

import (
	"github.com/labstack/echo"
)

type ITodos interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

type todos struct{}

// New Repository todos
func New() ITodos {
	return &todos{}
}

func (u todos) Create(c echo.Context) error {
	return nil
}
func (u todos) FindOne(c echo.Context) error {
	return nil
}
func (u todos) FindAll(c echo.Context) error {
	return nil
}
func (u todos) Update(c echo.Context) error {
	return nil
}
func (u todos) Delete(c echo.Context) error {
	return nil
}
func (u todos) Count(c echo.Context) error {
	return nil
}
