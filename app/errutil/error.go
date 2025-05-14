package errutil

import "errors"

var ResourceNotfound = errors.New("resource not found")
var RecordNotFound = errors.New("record not found")
var ActionNotFound = errors.New("action not found")
var AlreadyExists = errors.New("already exists")
var EmptyCollection = errors.New("empty collection")
