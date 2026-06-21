package service

import "errors"

var ErrInvalidStatusTransition = errors.New("cannot skip statuses")

var ErrInvalidPassword = errors.New("invalid password")

var ErrUserInactive = errors.New("user inactive")
var ErrUserAlreadyInactive = errors.New("user already inactive")

var ErrParcelNotDelivered = errors.New("cannot archive not delivered parcel")
var ErrParcelAlreadyArchived = errors.New("parcel already archive")
var ErrInvalidPage = errors.New("invalid page")
var ErrInvalidLimit = errors.New("invalid limit")
