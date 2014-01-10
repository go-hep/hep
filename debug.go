package rootio

import "fmt"

const g_rootio_debug = false

func myprintf(format string, args ...interface{}) (n int, err error) {
	if g_rootio_debug {
		return fmt.Printf(format, args...)
	}
	return
}
