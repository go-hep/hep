package rio

import (
	"errors"
)

var (
	ErrStreamNoRecMarker   = errors.New("rio: no record marker found")
	ErrRecordNoBlockMarker = errors.New("rio: no block marker found")
)

// EOF
