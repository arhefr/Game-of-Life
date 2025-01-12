package life

import "errors"

var (
	ErrPathFile    error = errors.New("error file path")
	ErrInvalidFile error = errors.New("error file data")

	ErrIncorrectArgs error = errors.New("error incorrect args")
)
