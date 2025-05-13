package backend

import "errors"

var (
	ErrorPath = errors.New("path must starts with /")
)
