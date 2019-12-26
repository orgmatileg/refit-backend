package todos

import (
	"github.com/labstack/echo"
)

// ITodos delivery http interface
type ITodos interface {
	Create(c echo.Context) error
	FindOneByID(c echo.Context) error
	FindAll(c echo.Context) error
	UpdateByID(c echo.Context) error
	DeleteByID(c echo.Context) error
}

type todos struct{}

// New delivery http repository todos
func New() ITodos {
	return &todos{}
}

func (u todos) Create(c echo.Context) error {
	return nil
}
func (u todos) FindOneByID(c echo.Context) error {
	return nil
}
func (u todos) FindAll(c echo.Context) error {
	return nil
}
func (u todos) UpdateByID(c echo.Context) error {
	return nil
}
func (u todos) DeleteByID(c echo.Context) error {
	return nil
}
