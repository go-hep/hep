package sio

import (
	"errors"
)

var (
	ErrStreamNoRecMarker   = errors.New("sio: no record marker found")
	ErrRecordNoBlockMarker = errors.New("sio: no block marker found")
	ErrBlockConnected      = errors.New("sio: block already connected")
)

// EOF
