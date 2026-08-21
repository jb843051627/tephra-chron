package model

import "errors"

var ErrMissing = errors.New("record missing")
var ErrInvalidState = errors.New("invalid lifecycle state")
var ErrRejected = errors.New("quality gate rejected")
