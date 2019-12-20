package helpers

import (
	"github.com/labstack/echo"
)

// ResponseDefault ...
type ResponseDefault struct {
	Code int           `json:"code"`
	Data []interface{} `json:"data"`
}

// MakeDefaultResponse ...
func MakeDefaultResponse(c echo.Context, httpStatusCode int, data ...interface{}) error {
	r := ResponseDefault{
		Code: httpStatusCode,
	}
	for _, d := range data {
		r.Data = append(r.Data, d)
	}
	return c.JSON(httpStatusCode, &r)
}
