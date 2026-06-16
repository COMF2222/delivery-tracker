package service

import "errors"

var ErrInvalidStatusTransition = errors.New("cannot skip statuses")
var ErrInvalidPassword = errors.New("invalid password")
var ErrUserInactive = errors.New("user inactive")
