package service

import "errors"

var ErrInvalidStatusTransition = errors.New("cannot skip statuses")
