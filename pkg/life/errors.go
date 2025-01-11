package life

import "errors"

var (
	ErrPathFile    error = errors.New("error file path")
	ErrInvalidFile error = errors.New("error file data")

	ErrWidthWorld      error = errors.New("error range: 100 <= width <= 600")
	ErrHeightWorld     error = errors.New("error range: 100 <= width <= 300")
	ErrChanceAliveCell error = errors.New("error range: 1 <= chance <= 100")
)
