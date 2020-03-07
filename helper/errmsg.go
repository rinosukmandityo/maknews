package helper

import (
	"errors"
)

var (
	ErrDataNotFound = errors.New("Data Not Found")
	ErrDataInvalid  = errors.New("Data Invalid")
)
