package helpers

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"regexp"
)

var (
	regexZeroTo100      = regexp.MustCompile("^(100|[1-9][0-9]?)$")
	regexPositiveNumber = regexp.MustCompile("^[+]?\\d+([.]\\d+)?$")
)

// ValidationQueryParamFindAll validation query param for find all endpoint
func ValidationQueryParamFindAll(limit, offset, order string) (err error) {
	return validation.Errors{
		"limit":  validation.Validate(limit, validation.Required, is.Digit, validation.Match(regexZeroTo100)),
		"offset": validation.Validate(offset, validation.Required, is.Digit, validation.Match(regexPositiveNumber)),
		"order":  validation.Validate(order, validation.Required, validation.In("asc", "desc")),
	}.Filter()
}
