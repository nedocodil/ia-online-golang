package controllers

import "errors"

var (
	ErrServer    = errors.New("error server")
	ErrParseJSON = errors.New("error parsing JSON")
)
