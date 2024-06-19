package utils

import "errors"

var (
	ErrIgnoreDirs = errors.New("ignoreDirs.json should be created in .godemon (It will be automatized later)")
)
