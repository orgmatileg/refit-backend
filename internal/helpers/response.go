package helpers

import (
	"fmt"
	"github.com/labstack/echo"
	"refit_backend/models"
	"reflect"
	"strconv"
)

// MakeDefaultResponse ...
func MakeDefaultResponse(c echo.Context, httpStatusCode int, data ...interface{}) error {
	r := models.ResponseDefault{
		Code: httpStatusCode,
	}

	for _, d := range data {

		if d == nil {
			r.Data = make([]interface{}, 0)
			return c.JSON(httpStatusCode, &r)
		}

		if reflect.TypeOf(d).Kind() == reflect.Slice {
			s := reflect.ValueOf(d)
			for i := 0; i < s.Len(); i++ {
				r.Data = append(r.Data, s.Index(i).Interface())
			}
		} else {
			switch v := d.(type) {
			case models.RequestQueryParam:
				backURL := v["back_url"]
				nextURL := v["next_url"]
				count := v["count"]
				limit := v["limit"]
				offset := v["offset"]
				order := v["order"]
				r.IsList(backURL, nextURL, count, limit, offset, order)
			case error:
				r.Data = append(r.Data, struct {
					Message string `json:"message"`
				}{v.Error()})
			default:
				r.Data = append(r.Data, d)
			}
		}

	}

	return c.JSON(httpStatusCode, &r)
}

// CreateBackAndNextURL for list endpoint
func CreateBackAndNextURL(c echo.Context, count uint, limit, offset, order string) (backURL, nextURL string, err error) {

	var (
		host = c.Request().Host
		path = c.Request().URL.Path
	)

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	currentPage := offsetInt / limitInt
	totalPage := int(count) / limitInt

	if offsetInt > 0 {
		offsetIntBackURL := offsetInt - limitInt
		backURL = fmt.Sprintf("%s%s?limit=%s&offset=%s&order=%s", host, path, limit, strconv.Itoa(offsetIntBackURL), order)
	}

	if totalPage > 1 && currentPage != totalPage {
		offsetInt += limitInt
		nextURL = fmt.Sprintf("%s%s?limit=%s&offset=%s&order=%s", host, path, limit, strconv.Itoa(offsetInt), order)
	} else {
		// return empty backURL and nextURL
		return backURL, nextURL, nil
	}

	return backURL, nextURL, nil
}
