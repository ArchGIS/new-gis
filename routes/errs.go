package routes

import (
	"errors"
)

var NotAllowedQueryParams = errors.New("Given query parameters in not allowed")
var NotValidQueryParameters = errors.New("Given parameters are invalid")
