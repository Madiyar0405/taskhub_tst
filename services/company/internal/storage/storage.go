package storage

import "errors"

var (
	ErrNoRecordFound = errors.New("no record found")
	AlreadyExists    = errors.New("already exists")
)
