package auth

import "errors"

var ErrUnauthorized error = errors.New("not authorized")
