package routes

import (
	"errors"
)

var NotAllowedQueryParams = errors.New("Given query parameters in not allowed")
