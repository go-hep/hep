package fwk

import (
	"fmt"
)

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// EOF
