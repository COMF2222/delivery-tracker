package repository

import "errors"

var ErrTrackNumberAlreadyExists = errors.New("track number already exists")
var ErrParcelNotFound = errors.New("parcel not found")
var ErrUserNotFound = errors.New("user not found")
