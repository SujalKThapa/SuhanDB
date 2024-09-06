package main

import "errors"

const (
	magicNumberSize = 4
	nodeHeaderSize  = 3
	pageNumSize     = 8
	collectionSize  = 16
	counterSize     = 4
)

var writeInsideReadTextErr = errors.New("Cannot perform a write transaction inside a read transaction.")
