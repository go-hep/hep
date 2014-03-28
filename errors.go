package rio

import (
	"errors"
)

var (
	ErrStreamNoRecMarker = errors.New("rio: no record marker found")
)

// EOF
