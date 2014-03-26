package fwk

import (
	"fmt"
)

func Errorf(format string, args ...interface{}) Error {
	return fmt.Errorf(format, args...)
}

// EOF
