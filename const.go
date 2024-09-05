package main

import "errors"

const (
	nodeHeaderSize = 3
	pageNumSize    = 8
)

var writeInsideReadTextErr = errors.New("Cannot perform a write transaction inside a read transaction.")
